const { execSync } = require('child_process');

/**
 * Get emoji indicator based on coverage percentage
 */
function getCoverageEmoji(percentage) {
  if (percentage >= 80) return '🟢';
  if (percentage >= 60) return '🟡';
  if (percentage >= 40) return '🟠';
  return '🔴';
}

/**
 * Parse coverage.out file and extract coverage data
 */
function parseCoverageFile(coverageFile) {
  // Get function-level coverage data
  const coverageOutput = execSync(`go tool cover -func=${coverageFile}`, { encoding: 'utf8' });

  // Extract total coverage
  const totalMatch = coverageOutput.match(/total:.*?(\d+\.\d+)%/);
  const totalCoverage = totalMatch ? `${totalMatch[1]}%` : '0.0%';

  // Parse per-package coverage
  const packageCoverage = {};
  const lines = coverageOutput.split('\n');

  for (const line of lines) {
    if (line.trim() && !line.includes('total:')) {
      const match = line.match(/^(.+?):\S+\s+(\d+\.\d+)%/);
      if (match) {
        const filePath = match[1];
        const coverage = parseFloat(match[2]);

        const parts = filePath.split('/');
        parts.pop();
        const packagePath = parts.join('/') || '.';

        if (!packageCoverage[packagePath]) {
          packageCoverage[packagePath] = { total: 0, count: 0 };
        }
        packageCoverage[packagePath].total += coverage;
        packageCoverage[packagePath].count += 1;
      }
    }
  }

  // Calculate and sort packages
  const packages = [];
  for (const [pkg, data] of Object.entries(packageCoverage)) {
    const avgCoverage = data.count > 0 ? data.total / data.count : 0;
    packages.push({ package: pkg, coverage: avgCoverage });
  }
  packages.sort((a, b) => b.coverage - a.coverage);

  return { totalCoverage, packages };
}

/**
 * Build the coverage comment markdown
 */
function buildCoverageComment(totalCoverage, packages, minThreshold) {
  const coverageNum = parseFloat(totalCoverage.replace('%', ''));
  const totalEmoji = getCoverageEmoji(coverageNum);

  // Build package table
  let packageTable = '| Package | Coverage |\n|---------|----------|\n';

  if (packages && packages.length > 0) {
    for (const { package: pkg, coverage } of packages) {
      const emoji = getCoverageEmoji(coverage);
      packageTable += `| \`${pkg}\` | ${emoji} ${coverage.toFixed(1)}% |\n`;
    }
  } else {
    packageTable += '| _No data_ | - |\n';
  }

  return `## 🧪 Go Test Coverage Report

### Overall Coverage
${totalEmoji} **${totalCoverage}** of code is covered by tests

---

### 📦 Coverage by Package
${packageTable}

---

<sub>Coverage threshold: ${minThreshold}% | 🟢 ≥80% 🟡 ≥60% 🟠 ≥40% 🔴 <40%</sub>`;
}

/**
 * Main function executed by github-script
 */
module.exports = async ({ github, context }) => {
  const coverageFile = process.env.COVERAGE_FILE;
  const minThreshold = parseFloat(process.env.MIN_THRESHOLD || '70.0');

  console.log(`Parsing coverage file: ${coverageFile}`);

  // Parse coverage file
  const { totalCoverage, packages } = parseCoverageFile(coverageFile);
  const coverageNum = parseFloat(totalCoverage.replace('%', ''));

  console.log(`Total coverage: ${totalCoverage}`);
  console.log(`Packages found: ${packages.length}`);

  // Check if we should skip commenting
  if (coverageNum > minThreshold) {
    console.log(`Coverage ${totalCoverage} is above threshold of ${minThreshold}%. Skipping comment.`);
    return;
  }

  // Build comment body
  const commentBody = buildCoverageComment(totalCoverage, packages, minThreshold);

  // Only comment on pull requests
  if (context.eventName !== 'pull_request') {
    console.log('Not a pull request event, skipping comment.');
    return;
  }

  // Find existing bot comment
  const { data: comments } = await github.rest.issues.listComments({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: context.issue.number,
  });

  const botComment = comments.find(
    comment => comment.user.type === 'Bot' &&
               comment.body.includes('## 🧪 Go Test Coverage Report')
  );

  if (botComment) {
    // Update existing comment
    await github.rest.issues.updateComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      comment_id: botComment.id,
      body: commentBody,
    });
    console.log(`Updated comment #${botComment.id}`);
  } else {
    // Create new comment
    await github.rest.issues.createComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      issue_number: context.issue.number,
      body: commentBody,
    });
    console.log('Created new coverage comment');
  }
};
