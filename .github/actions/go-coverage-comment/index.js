const core = require('@actions/core');
const github = require('@actions/github');
const fs = require('fs');
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
      // Format: path/to/file.go:function coverage%
      const match = line.match(/^(.+?):\S+\s+(\d+\.\d+)%/);
      if (match) {
        const filePath = match[1];
        const coverage = parseFloat(match[2]);

        // Extract package path (directory of the file)
        const parts = filePath.split('/');
        parts.pop(); // Remove filename
        const packagePath = parts.join('/') || '.';

        if (!packageCoverage[packagePath]) {
          packageCoverage[packagePath] = { total: 0, count: 0 };
        }
        packageCoverage[packagePath].total += coverage;
        packageCoverage[packagePath].count += 1;
      }
    }
  }

  // Calculate average coverage per package and format
  const packageResults = [];
  for (const [pkg, data] of Object.entries(packageCoverage)) {
    const avgCoverage = data.count > 0 ? data.total / data.count : 0;
    packageResults.push({ package: pkg, coverage: avgCoverage });
  }

  // Sort by coverage descending
  packageResults.sort((a, b) => b.coverage - a.coverage);

  return {
    totalCoverage,
    packages: packageResults
  };
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
 * Main action execution
 */
async function run() {
  try {
    // Get inputs
    const coverageFile = core.getInput('coverage-file', { required: true });
    const minThreshold = parseFloat(core.getInput('min-threshold', { required: false }) || '70.0');
    const token = core.getInput('github-token', { required: true });

    // Parse coverage file
    console.log(`Parsing coverage file: ${coverageFile}`);
    const { totalCoverage, packages } = parseCoverageFile(coverageFile);
    console.log(`Total coverage: ${totalCoverage}`);
    console.log(`Packages found: ${packages.length}`);

    const coverageNum = parseFloat(totalCoverage.replace('%', ''));

    // Check if we should skip commenting
    if (coverageNum > minThreshold) {
      console.log(`Coverage ${totalCoverage} is above threshold of ${minThreshold}%. Skipping comment.`);
      return;
    }

    // Build comment body
    const commentBody = buildCoverageComment(totalCoverage, packages, minThreshold);

    // Post or update comment
    const octokit = github.getOctokit(token);
    const context = github.context;

    if (context.eventName !== 'pull_request') {
      console.log('Not a pull request event, skipping comment.');
      return;
    }

    // Find existing bot comment
    const { data: comments } = await octokit.rest.issues.listComments({
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
      await octokit.rest.issues.updateComment({
        owner: context.repo.owner,
        repo: context.repo.repo,
        comment_id: botComment.id,
        body: commentBody,
      });
      console.log(`Updated comment #${botComment.id}`);
    } else {
      // Create new comment
      await octokit.rest.issues.createComment({
        owner: context.repo.owner,
        repo: context.repo.repo,
        issue_number: context.issue.number,
        body: commentBody,
      });
      console.log('Created new coverage comment');
    }

  } catch (error) {
    core.setFailed(error.message);
  }
}

run();
