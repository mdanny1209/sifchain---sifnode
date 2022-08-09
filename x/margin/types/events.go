//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

const EventOpen = "margin/mtp_open"
const EventClose = "margin/mtp_close"
const EventForceClose = "margin/mtp_force_close"
const EventIncrementalInterestPayment = "margin/mtp_incremental_interest_payment"
const EventInterestRateComputation = "margin/interest_rate_computation"
const EventMarginUpdateParams = "margin/update_params"
const EventRepayInsuranceFund = "margin/repay_insurance_fund"
const AttributeKeyPoolInterestRate = "margin_pool_interest_rate"
const AttributeKeyMarginParams = "margin_params"
