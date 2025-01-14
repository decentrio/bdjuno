CREATE OR REPLACE VIEW validator_summary AS
SELECT 
    vi.consensus_address AS address,
    COALESCE(vd.moniker, NULL) AS moniker,
    vvp.voting_power,
    COALESCE(vc.commission, NULL) AS commission,
    vs.status,
    vs.jailed
FROM 
    validator_info vi
LEFT JOIN 
    validator_voting_power vvp ON vi.consensus_address = vvp.validator_address
LEFT JOIN 
    validator_commission vc ON vi.consensus_address = vc.validator_address
LEFT JOIN 
    validator_status vs ON vi.consensus_address = vs.validator_address
LEFT JOIN 
    validator_description vd ON vi.consensus_address = vd.validator_address;

CREATE OR REPLACE VIEW account_summary AS
SELECT 
    ac.address AS address,
    COALESCE(bl.amount, NULL) AS amount,
    COALESCE(bl.denom, NULL) AS denom
FROM 
    account ac
LEFT JOIN 
    balance bl ON bl.address = ac.address;

    
