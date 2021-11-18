param (
	[Parameter(Mandatory, HelpMessage="Version number")]
	[string]$Version,

	[Parameter(Mandatory, HelpMessage="Commit message")]
	[string]$Message
)

go mod tidy
$testVal = go test ./...
if (-not($testVal.Contains('ok'))) {
	Write-Error "Module failed testing"
	Return
}

$testGit = git commit -a -m $Message
if ($testGit.Contains('nothing added')) {
	Write-Error "Failed to commit to git"
	Return
}
git tag $Version
$testGit = git push origin $Version

$testGo = go list -m github.com/mplecapt/SteamAchievementProgressGolang@$Version
if (-not($testGo.Contains('github.com/mplecapt/SteamAchievementProgressGolang '+$Version))) {
	Write-Error "Failed to list on go catalog"
	Return
}