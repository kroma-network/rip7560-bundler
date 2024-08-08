package config

import (
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/entities"
)

func NewReputationConstantsFromEnv() *entities.ReputationConstants {
	viper.SetDefault("rip7560_bundler_min_unstake_delay", 86400)
	viper.SetDefault("rip7560_bundler_min_stake_value", 2000000000000000)
	viper.SetDefault("rip7560_bundler_same_sender_mempool_count", 4)
	viper.SetDefault("rip7560_bundler_same_unstaked_entity_mempool_count", 11)
	viper.SetDefault("rip7560_bundler_throttled_entity_mempool_count", 4)
	viper.SetDefault("rip7560_bundler_throttled_entity_live_blocks", 10)
	viper.SetDefault("rip7560_bundler_throttled_entity_bundle_count", 4)
	viper.SetDefault("rip7560_bundler_min_inclusion_rate_denominator", 10)
	viper.SetDefault("rip7560_bundler_throttling_slack", 10)
	viper.SetDefault("rip7560_bundler_ban_slack", 50)

	_ = viper.BindEnv("rip7560_bundler_min_unstake_delay")
	_ = viper.BindEnv("rip7560_bundler_min_stake_value")
	_ = viper.BindEnv("rip7560_bundler_same_sender_mempool_count")
	_ = viper.BindEnv("rip7560_bundler_same_unstaked_entity_mempool_count")
	_ = viper.BindEnv("rip7560_bundler_throttled_entity_mempool_count")
	_ = viper.BindEnv("rip7560_bundler_throttled_entity_live_blocks")
	_ = viper.BindEnv("rip7560_bundler_throttled_entity_bundle_count")
	_ = viper.BindEnv("rip7560_bundler_min_inclusion_rate_denominator")
	_ = viper.BindEnv("rip7560_bundler_throttling_slack")
	_ = viper.BindEnv("rip7560_bundler_ban_slack")

	return &entities.ReputationConstants{
		MinUnstakeDelay:                viper.GetInt("rip7560_bundler_min_unstake_delay"),
		MinStakeValue:                  viper.GetInt64("rip7560_bundler_min_stake_value"),
		SameSenderMempoolCount:         viper.GetInt("rip7560_bundler_same_sender_mempool_count"),
		SameUnstakedEntityMempoolCount: viper.GetInt("rip7560_bundler_same_unstaked_entity_mempool_count"),
		ThrottledEntityMempoolCount:    viper.GetInt("rip7560_bundler_throttled_entity_mempool_count"),
		ThrottledEntityLiveBlocks:      viper.GetInt("rip7560_bundler_throttled_entity_live_blocks"),
		ThrottledEntityBundleCount:     viper.GetInt("rip7560_bundler_throttled_entity_bundle_count"),
		MinInclusionRateDenominator:    viper.GetInt("rip7560_bundler_min_inclusion_rate_denominator"),
		ThrottlingSlack:                viper.GetInt("rip7560_bundler_throttling_slack"),
		BanSlack:                       viper.GetInt("rip7560_bundler_ban_slack"),
	}
}
