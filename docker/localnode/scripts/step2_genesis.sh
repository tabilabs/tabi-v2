#!/usr/bin/env sh

# Input parameters
NODE_ID=${ID:-0}

echo "Preparing genesis file"

ACCOUNT_NAME="admin"
echo "Adding account $ACCOUNT_NAME"
printf "12345678\n12345678\ny\n" | tabid keys add $ACCOUNT_NAME >/dev/null 2>&1

override_genesis() {
  cat ~/.tabi/config/genesis.json | jq "$1" > ~/.tabi/config/tmp_genesis.json && mv ~/.tabi/config/tmp_genesis.json ~/.tabi/config/genesis.json;
}

override_genesis '.app_state["crisis"]["constant_fee"]["denom"]="atabi"'
override_genesis '.app_state["mint"]["params"]["mint_denom"]="atabi"'
override_genesis '.app_state["staking"]["params"]["bond_denom"]="atabi"'
override_genesis '.app_state["oracle"]["params"]["vote_period"]="2"'
override_genesis '.app_state["slashing"]["params"]["signed_blocks_window"]="10000"'
override_genesis '.app_state["slashing"]["params"]["min_signed_per_window"]="0.050000000000000000"'
override_genesis '.app_state["staking"]["params"]["max_validators"]="50"'
override_genesis '.consensus_params["block"]["max_gas"]="35000000"'
override_genesis '.app_state["staking"]["params"]["unbonding_time"]="10s"'

# Set a token release schedule for the genesis file
start_date="$(date +"%Y-%m-%d")"
end_date="$(date -d "+3 days" +"%Y-%m-%d")"
override_genesis ".app_state[\"mint\"][\"params\"][\"token_release_schedule\"]=[{\"start_date\": \"$start_date\", \"end_date\": \"$end_date\", \"token_release_amount\": \"999999999999\"}]"


# We already added node0's genesis account in configure_init, remove it here since we're going to re-add it in the "add genesis accounts" step
override_genesis '.app_state["auth"]["accounts"]=[]'
override_genesis '.app_state["bank"]["balances"]=[]'
override_genesis '.app_state["genutil"]["gen_txs"]=[]'
override_genesis '.app_state["bank"]["denom_metadata"]=[{"denom_units":[{"denom":"UATOM","exponent":6,"aliases":["UATOM"]}],"base":"uatom","display":"uatom","name":"UATOM","symbol":"UATOM"}]'

# gov parameters
override_genesis '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="atabi"'
override_genesis '.app_state["gov"]["deposit_params"]["min_expedited_deposit"][0]["denom"]="atabi"'
override_genesis '.app_state["gov"]["deposit_params"]["max_deposit_period"]="100s"'
override_genesis '.app_state["gov"]["voting_params"]["voting_period"]="30s"'
override_genesis '.app_state["gov"]["voting_params"]["expedited_voting_period"]="15s"'
override_genesis '.app_state["gov"]["tally_params"]["quorum"]="0.5"'
override_genesis '.app_state["gov"]["tally_params"]["threshold"]="0.5"'
override_genesis '.app_state["gov"]["tally_params"]["expedited_quorum"]="0.9"'
override_genesis '.app_state["gov"]["tally_params"]["expedited_threshold"]="0.9"'

# add genesis accounts for each node
while read account; do
  echo "Adding: $account"
  tabid add-genesis-account "$account" 1000000000000000000000atabi,1000000000000000000000uusdc,1000000000000000000000uatom
done <build/generated/genesis_accounts.txt

# add funds to admin account
printf "12345678\n" | tabid add-genesis-account admin 1000000000000000000000atabi,1000000000000000000000uusdc,1000000000000000000000uatom

mkdir -p ~/exported_keys
cp -r build/generated/gentx/* ~/.tabi/config/gentx
cp -r build/generated/exported_keys ~/exported_keys

# add validators to genesis
/usr/bin/add_validator_to_gensis.sh

# collect gentxs
echo "Collecting all gentx"
tabid collect-gentxs >/dev/null 2>&1

cp ~/.tabi/config/genesis.json build/generated/genesis.json
echo "Genesis file has been created successfully"
