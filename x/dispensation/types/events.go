package types

// dispensation module event types

const (
	AttributeValueCategory         = ModuleName
	EventTypeDistributionStarted   = "distribution_started"
	EventTypeDistributionRun       = "distribution_run"
	AttributeKeyFromModuleAccount  = "module_account"
	AttributeKeyDistributionName   = "distribution_name"
	AttributeKeyDistributionRunner = "distribution_runner"
	AttributeKeyDistributionList   = "distribution_list"
	AttributeKeyDistributionType   = "distribution_type"

	EventTypeClaimCreated = "userClaim_new"
	AttributeKeyClaimUser = "userClaim_creator"
	AttributeKeyClaimType = "userClaim_type"
	AttributeKeyClaimTime = "userClaim_creationTime"
)
