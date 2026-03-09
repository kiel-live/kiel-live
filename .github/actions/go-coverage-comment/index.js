const fs = require('fs');

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
 * Parse coverage.out file directly and extract coverage data
 * Format: mode: atomic
 *         path/to/file.go:startLine.startCol,endLine.endCol numStatements count
 */
function parseCoverageFile(coverageFile) {
  const content = fs.readFileSync(coverageFile, 'utf8');
  const lines = content.split('\n');

  const packageCoverage = {};
  let totalStatements = 0;
  let coveredStatements = 0;

  for (const line of lines) {
    // Skip mode line and empty lines
    if (!line.trim() || line.startsWith('mode:')) continue;

    // Parse: path/to/file.go:10.2,12.3 2 1
    const match = line.match(/^(.+?):[\d.,]+\s+(\d+)\s+(\d+)/);
    if (match) {
      const filePath = match[1];
      const numStatements = parseInt(match[2], 10);
      const count = parseInt(match[3], 10);

      // Extract package path (directory of the file)
      const parts = filePath.split('/');
      parts.pop(); // Remove filename
      const packagePath = parts.join('/') || '.';

      // Initialize package stats if needed
      if (!packageCoverage[packagePath]) {
        packageCoverage[packagePath] = { total: 0, covered: 0 };
      }

      // Update package coverage
      packageCoverage[packagePath].total += numStatements;
      if (count > 0) {
        packageCoverage[packagePath].covered += numStatements;
      }

      // Update total coverage
      totalStatements += numStatements;
      if (count > 0) {
        coveredStatements += numStatements;
      }
    }
  }

  // Calculate total coverage percentage
  const totalCoverage = totalStatements > 0
    ? `${((coveredStatements / totalStatements) * 100).toFixed(1)}%`
    : '0.0%';

  // Calculate per-package coverage and sort
  const packages = [];
  for (const [pkg, data] of Object.entries(packageCoverage)) {
    const coverage = data.total > 0 ? (data.covered / data.total) * 100 : 0;
    packages.push({ package: pkg, coverage });
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
