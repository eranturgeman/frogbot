//go:build integration

package main

import (
	"github.com/jfrog/frogbot/v3/utils"
	"github.com/jfrog/froggit-go/vcsclient"
	"github.com/jfrog/froggit-go/vcsutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	//#nosec G101 -- False positive - no hardcoded credentials.
	bitbucketCloudIntegrationTokenEnv = "FROGBOT_V3_TESTS_BITBUCKET_CLOUD_TOKEN"
	bitbucketCloudGitCloneUrl         = "https://bitbucket.org/frogbot-e2e-test/frogbot-test.git"
	bitbucketCloudRepoOwner           = "frogbot-e2e-test"
	// Atlassian API tokens with Bitbucket scopes require the literal username "x-bitbucket-api-token-auth"
	// for git's HTTP transport (clone/push), but that username is rejected by Bitbucket Cloud's REST API,
	// which instead expects Bearer auth (no username at all). The two credentials can't be unified, so the
	// git-only username is kept separate from the (empty) REST/env-var username via GitPushUsername.
	bitbucketCloudGitPushUsername = "x-bitbucket-api-token-auth"
)

func buildBitbucketCloudClient(t *testing.T, bitbucketCloudToken string) vcsclient.VcsClient {
	bbClient, err := vcsclient.NewClientBuilder(vcsutils.BitbucketCloud).Token(bitbucketCloudToken).Build()
	assert.NoError(t, err)
	return bbClient
}

func buildBitbucketCloudIntegrationTestDetails(t *testing.T) *IntegrationTestDetails {
	integrationRepoToken := getIntegrationToken(t, bitbucketCloudIntegrationTokenEnv)
	testDetails := NewIntegrationTestDetails(integrationRepoToken, string(utils.BitbucketCloud), bitbucketCloudGitCloneUrl, bitbucketCloudRepoOwner)
	testDetails.GitUsername = ""
	testDetails.GitPushUsername = bitbucketCloudGitPushUsername
	return testDetails
}

func bitbucketCloudTestsInit(t *testing.T) (vcsclient.VcsClient, *IntegrationTestDetails) {
	testDetails := buildBitbucketCloudIntegrationTestDetails(t)
	bbClient := buildBitbucketCloudClient(t, testDetails.GitToken)
	return bbClient, testDetails
}

func TestBitbucketCloud_ScanPullRequestIntegration(t *testing.T) {
	bbClient, testDetails := bitbucketCloudTestsInit(t)
	runScanPullRequestCmd(t, bbClient, testDetails)
}

func TestBitbucketCloud_ScanRepositoryIntegration(t *testing.T) {
	bbClient, testDetails := bitbucketCloudTestsInit(t)
	runScanRepositoryCmd(t, bbClient, testDetails)
}

func TestHelper_CleanupIntegrationTestsArtifactsBitbucketCloud(t *testing.T) {
	bbClient, testDetails := bitbucketCloudTestsInit(t)
	cleanupIntegrationArtifacts(t, bbClient, testDetails)
}
