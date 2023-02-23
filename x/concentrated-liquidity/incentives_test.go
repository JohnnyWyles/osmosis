package concentrated_liquidity_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cl "github.com/osmosis-labs/osmosis/v14/x/concentrated-liquidity"
	"github.com/osmosis-labs/osmosis/v14/x/concentrated-liquidity/model"
	"github.com/osmosis-labs/osmosis/v14/x/concentrated-liquidity/types"
)

var (
	defaultPoolId = uint64(1)

	testAddressOne   = sdk.AccAddress([]byte("addr1_______________"))
	testAddressTwo   = sdk.AccAddress([]byte("addr2_______________"))
	testAddressThree = sdk.AccAddress([]byte("addr3_______________"))

	testAccumOne = "testAccumOne"

	// Note: lexicographic order is denomFour, denomOne, denomThree, denomTwo
	testDenomOne   = "denomOne"
	testDenomTwo   = "denomTwo"
	testDenomThree = "denomThree"
	testDenomFour  = "denomFour"

	defaultIncentiveAmount = sdk.NewDec(2 << 60)

	testEmissionOne   = sdk.MustNewDecFromStr("0.000001")
	testEmissionTwo   = sdk.MustNewDecFromStr("0.0783")
	testEmissionThree = sdk.MustNewDecFromStr("165.4")
	testEmissionFour  = sdk.MustNewDecFromStr("57.93")

	defaultBlockTime  = time.Unix(0, 0).UTC()
	defaultTimeBuffer = time.Hour
	defaultStartTime  = defaultBlockTime.Add(defaultTimeBuffer)

	testUptimeOne   = types.SupportedUptimes[0]
	testUptimeTwo   = types.SupportedUptimes[1]
	testUptimeThree = types.SupportedUptimes[2]
	testUptimeFour  = types.SupportedUptimes[3]

	incentiveRecordOne = types.IncentiveRecord{
		PoolId:          validPoolId,
		IncentiveDenom:  testDenomOne,
		RemainingAmount: defaultIncentiveAmount,
		EmissionRate:    testEmissionOne,
		StartTime:       defaultStartTime,
		MinUptime:       testUptimeOne,
	}

	incentiveRecordTwo = types.IncentiveRecord{
		PoolId:          validPoolId,
		IncentiveDenom:  testDenomTwo,
		RemainingAmount: defaultIncentiveAmount,
		EmissionRate:    testEmissionTwo,
		StartTime:       defaultStartTime,
		MinUptime:       testUptimeTwo,
	}

	incentiveRecordThree = types.IncentiveRecord{
		PoolId:          validPoolId,
		IncentiveDenom:  testDenomThree,
		RemainingAmount: defaultIncentiveAmount,
		EmissionRate:    testEmissionThree,
		StartTime:       defaultStartTime,
		MinUptime:       testUptimeThree,
	}

	incentiveRecordFour = types.IncentiveRecord{
		PoolId:          validPoolId,
		IncentiveDenom:  testDenomFour,
		RemainingAmount: defaultIncentiveAmount,
		EmissionRate:    testEmissionFour,
		StartTime:       defaultStartTime,
		MinUptime:       testUptimeFour,
	}

	testQualifyingDepositsOne   = sdk.NewInt(50)
	testQualifyingDepositsTwo   = sdk.NewInt(100)
	testQualifyingDepositsThree = sdk.NewInt(399)
)

type ExpectedUptimes struct {
	emptyExpectedAccumValues     []sdk.DecCoins
	hundredTokensSingleDenom     []sdk.DecCoins
	hundredTokensMultiDenom      []sdk.DecCoins
	twoHundredTokensMultiDenom   []sdk.DecCoins
	threeHundredTokensMultiDenom []sdk.DecCoins
	fourHundredTokensMultiDenom  []sdk.DecCoins
	varyingTokensSingleDenom     []sdk.DecCoins
	varyingTokensMultiDenom      []sdk.DecCoins
}

// getExpectedUptimes returns a base set of expected values for testing based on the number
// of supported uptimes at runtime. This abstraction exists only to ensure backwards-compatibility
// of incentives-related tests if the supported uptimes are ever changed.
func getExpectedUptimes() ExpectedUptimes {
	expUptimes := ExpectedUptimes{
		emptyExpectedAccumValues:     []sdk.DecCoins{},
		hundredTokensSingleDenom:     []sdk.DecCoins{},
		hundredTokensMultiDenom:      []sdk.DecCoins{},
		twoHundredTokensMultiDenom:   []sdk.DecCoins{},
		threeHundredTokensMultiDenom: []sdk.DecCoins{},
		fourHundredTokensMultiDenom:  []sdk.DecCoins{},
		varyingTokensSingleDenom:     []sdk.DecCoins{},
		varyingTokensMultiDenom:      []sdk.DecCoins{},
	}
	for i := range types.SupportedUptimes {
		expUptimes.emptyExpectedAccumValues = append(expUptimes.emptyExpectedAccumValues, cl.EmptyCoins)
		expUptimes.hundredTokensSingleDenom = append(expUptimes.hundredTokensSingleDenom, sdk.NewDecCoins(cl.HundredFooCoins))
		expUptimes.hundredTokensMultiDenom = append(expUptimes.hundredTokensMultiDenom, sdk.NewDecCoins(cl.HundredFooCoins, cl.HundredBarCoins))
		expUptimes.twoHundredTokensMultiDenom = append(expUptimes.twoHundredTokensMultiDenom, sdk.NewDecCoins(cl.HundredFooCoins.Add(cl.HundredFooCoins), cl.HundredBarCoins.Add(cl.HundredBarCoins)))
		expUptimes.threeHundredTokensMultiDenom = append(expUptimes.threeHundredTokensMultiDenom, sdk.NewDecCoins(cl.HundredFooCoins.Add(cl.HundredFooCoins).Add(cl.HundredFooCoins), cl.HundredBarCoins.Add(cl.HundredBarCoins).Add(cl.HundredBarCoins)))
		expUptimes.fourHundredTokensMultiDenom = append(expUptimes.fourHundredTokensMultiDenom, sdk.NewDecCoins(cl.HundredFooCoins.Add(cl.HundredFooCoins).Add(cl.HundredFooCoins).Add(cl.HundredFooCoins), cl.HundredBarCoins.Add(cl.HundredBarCoins).Add(cl.HundredBarCoins).Add(cl.HundredBarCoins)))
		expUptimes.varyingTokensSingleDenom = append(expUptimes.varyingTokensSingleDenom, sdk.NewDecCoins(cl.HundredFooCoins.Add(sdk.NewDecCoin("foo", sdk.NewInt(int64(i))))))
		expUptimes.varyingTokensMultiDenom = append(expUptimes.varyingTokensMultiDenom, sdk.NewDecCoins(cl.HundredFooCoins.Add(sdk.NewDecCoin("foo", sdk.NewInt(int64(i)))), cl.HundredBarCoins.Add(sdk.NewDecCoin("bar", sdk.NewInt(int64(i*3))))))
	}

	return expUptimes
}

// Helper for converting raw DecCoins accum values to pool proto compatible UptimeTrackers
func wrapUptimeTrackers(accumValues []sdk.DecCoins) []model.UptimeTracker {
	wrappedUptimeTrackers := []model.UptimeTracker{}
	for _, accumValue := range accumValues {
		wrappedUptimeTrackers = append(wrappedUptimeTrackers, model.UptimeTracker{accumValue})
	}

	return wrappedUptimeTrackers
}

func expectedIncentives(denom string, rate sdk.Dec, timeElapsed time.Duration, qualifyingLiquidity sdk.Dec) sdk.DecCoin {
	timeInSec := sdk.NewDec(int64(timeElapsed)).Quo(sdk.MustNewDecFromStr("1000000000"))
	amount := rate.Mul(timeInSec).QuoTruncate(qualifyingLiquidity)

	return sdk.NewDecCoinFromDec(denom, amount)
}

func chargeIncentive(incentiveRecord types.IncentiveRecord, timeElapsed time.Duration) types.IncentiveRecord {
	incentivesEmitted := incentiveRecord.EmissionRate.Mul(sdk.NewDec(int64(timeElapsed)).Quo(sdk.MustNewDecFromStr("1000000000")))
	incentiveRecord.RemainingAmount = incentiveRecord.RemainingAmount.Sub(incentivesEmitted)

	return incentiveRecord
}

// Helper for adding a predetermined amount to each global uptime accum in clPool
func addToUptimeAccums(ctx sdk.Context, poolId uint64, clKeeper *cl.Keeper, addValues []sdk.DecCoins) error {
	poolUptimeAccumulators, err := clKeeper.GetUptimeAccumulators(ctx, poolId)
	if err != nil {
		return err
	}

	for uptimeIndex, uptimeAccum := range poolUptimeAccumulators {
		uptimeAccum.AddToAccumulator(addValues[uptimeIndex])
	}

	return nil
}

func (s *KeeperTestSuite) TestCreateAndGetUptimeAccumulators() {
	// We expect there to be len(types.SupportedUptimes) number of initialized accumulators
	// for a successful pool creation. We calculate this upfront to ensure test compatibility
	// if the uptimes we support ever change.
	expectedUptimes := getExpectedUptimes()

	type initUptimeAccumTest struct {
		poolId              uint64
		initializePoolAccum bool
		expectedAccumValues []sdk.DecCoins

		expectedPass bool
	}
	tests := map[string]initUptimeAccumTest{
		"default pool setup": {
			poolId:              defaultPoolId,
			initializePoolAccum: true,
			expectedAccumValues: expectedUptimes.emptyExpectedAccumValues,
			expectedPass:        true,
		},
		"setup with different poolId": {
			poolId:              defaultPoolId + 1,
			initializePoolAccum: true,
			expectedAccumValues: expectedUptimes.emptyExpectedAccumValues,
			expectedPass:        true,
		},
		"pool not initialized": {
			initializePoolAccum: false,
			poolId:              defaultPoolId,
			expectedAccumValues: []sdk.DecCoins{},
			expectedPass:        false,
		},
	}

	for name, tc := range tests {
		tc := tc
		s.Run(name, func() {
			s.SetupTest()
			clKeeper := s.App.ConcentratedLiquidityKeeper

			// system under test
			if tc.initializePoolAccum {
				err := clKeeper.CreateUptimeAccumulators(s.Ctx, tc.poolId)
				s.Require().NoError(err)
			}
			poolUptimeAccumulators, err := clKeeper.GetUptimeAccumulators(s.Ctx, tc.poolId)

			if tc.expectedPass {
				s.Require().NoError(err)

				// ensure number of uptime accumulators match supported uptimes
				s.Require().Equal(len(tc.expectedAccumValues), len(poolUptimeAccumulators))

				// ensure that each uptime was initialized to the correct value (sdk.DecCoins(nil))
				accumValues := []sdk.DecCoins{}
				for _, accum := range poolUptimeAccumulators {
					accumValues = append(accumValues, accum.GetValue())
				}
				s.Require().Equal(tc.expectedAccumValues, accumValues)
			} else {
				s.Require().Error(err)

				// ensure no accumulators exist for an uninitialized pool
				s.Require().Equal(0, len(poolUptimeAccumulators))
			}
		})
	}
}

func (s *KeeperTestSuite) TestGetUptimeAccumulatorName() {
	type getUptimeNameTest struct {
		poolId            uint64
		uptimeIndex       uint64
		expectedAccumName string
	}
	tests := map[string]getUptimeNameTest{
		"pool id 1, uptime id 0": {
			poolId:            defaultPoolId,
			uptimeIndex:       uint64(0),
			expectedAccumName: "uptime/1/0",
		},
		"pool id 1, uptime id 999": {
			poolId:            defaultPoolId,
			uptimeIndex:       uint64(999),
			expectedAccumName: "uptime/1/999",
		},
		"pool id 999, uptime id 1": {
			poolId:            uint64(999),
			uptimeIndex:       uint64(1),
			expectedAccumName: "uptime/999/1",
		},
	}

	for name, tc := range tests {
		tc := tc
		s.Run(name, func() {
			s.SetupTest()

			// system under test
			accumName := cl.GetUptimeAccumulatorName(tc.poolId, tc.uptimeIndex)
			s.Require().Equal(tc.expectedAccumName, accumName)
		})
	}
}

func (s *KeeperTestSuite) TestCreateAndGetUptimeAccumulatorValues() {
	// We expect there to be len(types.SupportedUptimes) number of initialized accumulators
	// for a successful pool creation.
	// We re-calculate these values each time to ensure test compatibility if the uptimes
	// we support ever change.
	expectedUptimes := getExpectedUptimes()

	type initUptimeAccumTest struct {
		poolId               uint64
		initializePoolAccums bool
		addedAccumValues     []sdk.DecCoins
		numTimesAdded        int
		expectedAccumValues  []sdk.DecCoins

		expectedPass bool
	}
	tests := map[string]initUptimeAccumTest{
		"hundred of a single denom in each accumulator added once": {
			poolId:               defaultPoolId,
			initializePoolAccums: true,
			addedAccumValues:     expectedUptimes.hundredTokensSingleDenom,
			numTimesAdded:        1,
			expectedAccumValues:  expectedUptimes.hundredTokensSingleDenom,
			expectedPass:         true,
		},
		"hundred of multiple denom in each accumulator added once": {
			poolId:               defaultPoolId,
			initializePoolAccums: true,
			addedAccumValues:     expectedUptimes.hundredTokensMultiDenom,
			numTimesAdded:        1,
			expectedAccumValues:  expectedUptimes.hundredTokensMultiDenom,
			expectedPass:         true,
		},
		"varying amounts of single denom in each accumulator added once": {
			poolId:               defaultPoolId,
			initializePoolAccums: true,
			addedAccumValues:     expectedUptimes.varyingTokensSingleDenom,
			numTimesAdded:        1,
			expectedAccumValues:  expectedUptimes.varyingTokensSingleDenom,
			expectedPass:         true,
		},
		"varying of multiple denoms in each accumulator added once": {
			poolId:               defaultPoolId,
			initializePoolAccums: true,
			addedAccumValues:     expectedUptimes.varyingTokensMultiDenom,
			numTimesAdded:        1,
			expectedAccumValues:  expectedUptimes.varyingTokensMultiDenom,
			expectedPass:         true,
		},
		"hundred of multiple denom in each accumulator added twice": {
			poolId:               defaultPoolId,
			initializePoolAccums: true,
			addedAccumValues:     expectedUptimes.hundredTokensMultiDenom,
			numTimesAdded:        2,
			expectedAccumValues:  expectedUptimes.twoHundredTokensMultiDenom,
			expectedPass:         true,
		},
		"setup with different poolId": {
			poolId:               defaultPoolId + 1,
			initializePoolAccums: true,
			addedAccumValues:     expectedUptimes.hundredTokensSingleDenom,
			numTimesAdded:        1,
			expectedAccumValues:  expectedUptimes.hundredTokensSingleDenom,
			expectedPass:         true,
		},
		"pool not initialized": {
			initializePoolAccums: false,
			poolId:               defaultPoolId,
			addedAccumValues:     expectedUptimes.hundredTokensSingleDenom,
			numTimesAdded:        1,
			expectedAccumValues:  []sdk.DecCoins{},
			expectedPass:         false,
		},
	}

	for name, tc := range tests {
		tc := tc
		s.Run(name, func() {
			s.SetupTest()
			clKeeper := s.App.ConcentratedLiquidityKeeper

			// system under test
			var err error
			if tc.initializePoolAccums {
				err = clKeeper.CreateUptimeAccumulators(s.Ctx, tc.poolId)
				s.Require().NoError(err)

				poolUptimeAccumulators, err := clKeeper.GetUptimeAccumulators(s.Ctx, tc.poolId)
				s.Require().NoError(err)

				for i := 0; i < tc.numTimesAdded; i++ {
					for uptimeId, uptimeAccum := range poolUptimeAccumulators {
						uptimeAccum.AddToAccumulator(tc.addedAccumValues[uptimeId])
					}
					poolUptimeAccumulators, err = clKeeper.GetUptimeAccumulators(s.Ctx, tc.poolId)
					s.Require().NoError(err)
				}
			}
			poolUptimeAccumulatorValues, err := clKeeper.GetUptimeAccumulatorValues(s.Ctx, tc.poolId)

			if tc.expectedPass {
				s.Require().NoError(err)

				// ensure number of uptime accumulators match supported uptimes
				s.Require().Equal(len(tc.expectedAccumValues), len(poolUptimeAccumulatorValues))

				// ensure that each uptime was initialized to the correct value (sdk.DecCoins(nil))
				s.Require().Equal(tc.expectedAccumValues, poolUptimeAccumulatorValues)
			} else {
				s.Require().Error(err)

				// ensure no accumulators exist for an uninitialized pool
				s.Require().Equal(0, len(poolUptimeAccumulatorValues))
			}
		})
	}
}

func (s *KeeperTestSuite) TestCalcAccruedIncentivesForAccum() {
	type calcAccruedIncentivesTest struct {
		poolId               uint64
		accumUptime          time.Duration
		qualifyingLiquidity  sdk.Dec
		timeElapsed          time.Duration
		poolIncentiveRecords []types.IncentiveRecord

		expectedResult           sdk.DecCoins
		expectedIncentiveRecords []types.IncentiveRecord
		expectedPass             bool
	}
	tests := map[string]calcAccruedIncentivesTest{
		"one incentive record, one qualifying for incentives": {
			poolId:               defaultPoolId,
			accumUptime:          types.SupportedUptimes[0],
			qualifyingLiquidity:  sdk.NewDec(100),
			timeElapsed:          time.Hour,
			poolIncentiveRecords: []types.IncentiveRecord{incentiveRecordOne},

			expectedResult: sdk.DecCoins{
				expectedIncentives(incentiveRecordOne.IncentiveDenom, incentiveRecordOne.EmissionRate, time.Hour, sdk.NewDec(100)),
			},
			expectedIncentiveRecords: []types.IncentiveRecord{chargeIncentive(incentiveRecordOne, time.Hour)},
			expectedPass:             true,
		},
		"two incentive records, one qualifying for incentives": {
			poolId:               defaultPoolId,
			accumUptime:          types.SupportedUptimes[0],
			qualifyingLiquidity:  sdk.NewDec(100),
			timeElapsed:          time.Hour,
			poolIncentiveRecords: []types.IncentiveRecord{incentiveRecordOne, incentiveRecordTwo},

			expectedResult: sdk.DecCoins{
				// We only expect the first incentive record to qualify
				expectedIncentives(incentiveRecordOne.IncentiveDenom, incentiveRecordOne.EmissionRate, time.Hour, sdk.NewDec(100)),
			},
			expectedIncentiveRecords: []types.IncentiveRecord{
				// We only charge the first incentive record since the second wasn't affected
				chargeIncentive(incentiveRecordOne, time.Hour),
				incentiveRecordTwo,
			},
			expectedPass: true,
		},

		// error catching
		"zero qualifying liquidity": {
			poolId:               defaultPoolId,
			accumUptime:          types.SupportedUptimes[0],
			qualifyingLiquidity:  sdk.NewDec(0),
			timeElapsed:          time.Hour,
			poolIncentiveRecords: []types.IncentiveRecord{incentiveRecordOne},

			expectedResult:           sdk.DecCoins{},
			expectedIncentiveRecords: []types.IncentiveRecord{},
			expectedPass:             false,
		},
		"zero time elapsed": {
			poolId:               defaultPoolId,
			accumUptime:          types.SupportedUptimes[0],
			qualifyingLiquidity:  sdk.NewDec(100),
			timeElapsed:          time.Duration(0),
			poolIncentiveRecords: []types.IncentiveRecord{incentiveRecordOne},

			expectedResult:           sdk.DecCoins{},
			expectedIncentiveRecords: []types.IncentiveRecord{},
			expectedPass:             false,
		},
	}

	for name, tc := range tests {
		tc := tc
		s.Run(name, func() {
			s.SetupTest()
			s.Ctx = s.Ctx.WithBlockTime(defaultStartTime.Add(tc.timeElapsed))

			s.PrepareConcentratedPool()

			// system under test
			actualResult, updatedPoolRecords, err := cl.CalcAccruedIncentivesForAccum(s.Ctx, tc.accumUptime, tc.qualifyingLiquidity, sdk.NewDec(int64(tc.timeElapsed)).Quo(sdk.MustNewDecFromStr("1000000000")), tc.poolIncentiveRecords)

			if tc.expectedPass {
				s.Require().NoError(err)

				s.Require().Equal(tc.expectedResult, actualResult)
				s.Require().Equal(tc.expectedIncentiveRecords, updatedPoolRecords)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateUptimeAccumulatorsToNow() {
	supportedUptimes := types.SupportedUptimes

	type updateAccumToNow struct {
		poolId               uint64
		accumUptime          time.Duration
		qualifyingLiquidity  sdk.Dec
		timeElapsed          time.Duration
		poolIncentiveRecords []types.IncentiveRecord

		expectedResult           sdk.DecCoins
		expectedUptimeDeltas     []sdk.DecCoins
		expectedIncentiveRecords []types.IncentiveRecord
		expectedPass             bool
	}
	tests := map[string]updateAccumToNow{
		"one incentive record": {
			poolId:               defaultPoolId,
			timeElapsed:          time.Hour,
			poolIncentiveRecords: []types.IncentiveRecord{incentiveRecordOne},

			expectedIncentiveRecords: []types.IncentiveRecord{
				// We deduct incentives from the record for the period it emitted incentives
				chargeIncentive(incentiveRecordOne, time.Hour),
			},
			expectedPass: true,
		},
		"two incentive records, each with qualifying liquidity": {
			poolId:               defaultPoolId,
			timeElapsed:          time.Hour,
			poolIncentiveRecords: []types.IncentiveRecord{incentiveRecordOne, incentiveRecordTwo},

			expectedIncentiveRecords: []types.IncentiveRecord{
				// We deduct incentives from both records since there are positions for each
				chargeIncentive(incentiveRecordOne, time.Hour),
				chargeIncentive(incentiveRecordTwo, time.Hour),
			},
			expectedPass: true,
		},
		"three incentive records, each with qualifying liquidity": {
			poolId:               defaultPoolId,
			timeElapsed:          time.Hour,
			poolIncentiveRecords: []types.IncentiveRecord{incentiveRecordOne, incentiveRecordTwo, incentiveRecordThree},

			expectedIncentiveRecords: []types.IncentiveRecord{
				// We deduct incentives from each record since there are positions for all three
				// Note that records are ordered lexicographically by denom in state
				chargeIncentive(incentiveRecordOne, time.Hour),
				chargeIncentive(incentiveRecordThree, time.Hour),
				chargeIncentive(incentiveRecordTwo, time.Hour),
			},
			expectedPass: true,
		},
		"two incentive records, only one with qualifying liquidity": {
			poolId:               defaultPoolId,
			timeElapsed:          time.Hour,
			poolIncentiveRecords: []types.IncentiveRecord{incentiveRecordOne, incentiveRecordTwo, incentiveRecordThree, incentiveRecordFour},

			expectedIncentiveRecords: []types.IncentiveRecord{
				// We only deduct from the first three incentive records since the last doesn't emit anything
				// Note that records are ordered lexicographically by denom in state
				incentiveRecordFour,
				chargeIncentive(incentiveRecordOne, time.Hour),
				chargeIncentive(incentiveRecordThree, time.Hour),
				chargeIncentive(incentiveRecordTwo, time.Hour),
			},
			expectedPass: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		s.Run(name, func() {
			s.SetupTest()
			clKeeper := s.App.ConcentratedLiquidityKeeper
			s.Ctx = s.Ctx.WithBlockTime(defaultStartTime)

			// Set up test pool
			clPool := s.PrepareConcentratedPool()

			// Initialize test incentives on the pool
			clKeeper.SetMultipleIncentiveRecords(s.Ctx, tc.poolIncentiveRecords)

			// Get initial uptime accum values for comparison
			initUptimeAccumValues, err := clKeeper.GetUptimeAccumulatorValues(s.Ctx, tc.poolId)
			s.Require().NoError(err)

			// Add qualifying and non-qualifying liquidity to the pool
			s.FundAcc(testAddressOne, sdk.NewCoins(sdk.NewCoin(clPool.GetToken0(), testQualifyingDepositsOne), sdk.NewCoin(clPool.GetToken1(), testQualifyingDepositsOne)))
			s.FundAcc(testAddressTwo, sdk.NewCoins(sdk.NewCoin(clPool.GetToken0(), testQualifyingDepositsTwo), sdk.NewCoin(clPool.GetToken1(), testQualifyingDepositsTwo)))
			s.FundAcc(testAddressThree, sdk.NewCoins(sdk.NewCoin(clPool.GetToken0(), testQualifyingDepositsThree), sdk.NewCoin(clPool.GetToken1(), testQualifyingDepositsThree)))

			_, _, qualifyingLiquidityUptimeOne, err := clKeeper.CreatePosition(s.Ctx, tc.poolId, testAddressOne, testQualifyingDepositsOne, testQualifyingDepositsOne, sdk.ZeroInt(), sdk.ZeroInt(), clPool.GetCurrentTick().Int64()-1, clPool.GetCurrentTick().Int64()+1, defaultStartTime.Add(supportedUptimes[0]))
			s.Require().NoError(err)

			_, _, qualifyingLiquidityUptimeTwo, err := clKeeper.CreatePosition(s.Ctx, tc.poolId, testAddressTwo, testQualifyingDepositsTwo, testQualifyingDepositsTwo, sdk.ZeroInt(), sdk.ZeroInt(), clPool.GetCurrentTick().Int64()-1, clPool.GetCurrentTick().Int64()+1, defaultStartTime.Add(supportedUptimes[1]))
			s.Require().NoError(err)

			_, _, qualifyingLiquidityUptimeThree, err := clKeeper.CreatePosition(s.Ctx, tc.poolId, testAddressThree, testQualifyingDepositsThree, testQualifyingDepositsThree, sdk.ZeroInt(), sdk.ZeroInt(), clPool.GetCurrentTick().Int64()-1, clPool.GetCurrentTick().Int64()+1, defaultStartTime.Add(supportedUptimes[2]))
			s.Require().NoError(err)

			// Note that the third position (1D freeze) qualifies for all three uptimes, the second position qualifies for the first two,
			// and the first position only qualifies for the first. None of the positions qualify for any later uptimes (e.g. 1W)
			qualifyingLiquidities := []sdk.Dec{
				qualifyingLiquidityUptimeOne.Add(qualifyingLiquidityUptimeTwo).Add(qualifyingLiquidityUptimeThree),
				qualifyingLiquidityUptimeTwo.Add(qualifyingLiquidityUptimeThree),
				qualifyingLiquidityUptimeThree,
			}

			// Let `timeElapsed` time pass
			s.Ctx = s.Ctx.WithBlockTime(defaultStartTime.Add(tc.timeElapsed))

			// system under test
			err = clKeeper.UpdateUptimeAccumulatorsToNow(s.Ctx, tc.poolId)

			if tc.expectedPass {
				s.Require().NoError(err)

				// Get updated pool for testing purposes
				clPool, err := clKeeper.GetPoolById(s.Ctx, tc.poolId)
				s.Require().NoError(err)

				// Get new uptime accum values for comparison
				newUptimeAccumValues, err := clKeeper.GetUptimeAccumulatorValues(s.Ctx, tc.poolId)
				s.Require().NoError(err)

				// Calculate expected uptime deltas using qualifying liquidity deltas (eh can only test one incentive?)
				expectedUptimeDeltas := []sdk.DecCoins{}
				for uptimeIndex := range newUptimeAccumValues {
					if uptimeIndex < len(tc.poolIncentiveRecords) && uptimeIndex < len(qualifyingLiquidities) {
						expectedUptimeDeltas = append(expectedUptimeDeltas, sdk.NewDecCoins(expectedIncentives(tc.poolIncentiveRecords[uptimeIndex].IncentiveDenom, tc.poolIncentiveRecords[uptimeIndex].EmissionRate, time.Hour, qualifyingLiquidities[uptimeIndex])))
					} else {
						expectedUptimeDeltas = append(expectedUptimeDeltas, cl.EmptyCoins)
					}
				}

				// Ensure that each accumulator value changes by the correct amount
				for uptimeIndex := range newUptimeAccumValues {
					uptimeDelta := newUptimeAccumValues[uptimeIndex].Sub(initUptimeAccumValues[uptimeIndex])
					s.Require().Equal(expectedUptimeDeltas[uptimeIndex], uptimeDelta)
				}

				// Ensure that LastLiquidityUpdate field is updated for pool
				s.Require().Equal(s.Ctx.BlockTime(), clPool.GetLastLiquidityUpdate())
				// Ensure that pool's IncentiveRecords are updated to reflect emitted incentives
				updatedIncentiveRecords, err := clKeeper.GetAllIncentiveRecordsForPool(s.Ctx, tc.poolId)
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedIncentiveRecords, updatedIncentiveRecords)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

// Note: we test that incentive records are properly deducted by emissions in `TestUpdateUptimeAccumulatorsToNow` above.
// This test aims to cover the behavior of a series of state read/writes relating to incentive records.
func (s *KeeperTestSuite) TestIncentiveRecordsSetAndGet() {
	s.SetupTest()
	clKeeper := s.App.ConcentratedLiquidityKeeper
	s.Ctx = s.Ctx.WithBlockTime(defaultStartTime)
	emptyIncentiveRecords := []types.IncentiveRecord{}

	// Set up test pool
	clPoolOne := s.PrepareConcentratedPool()

	// Set up second pool for reference
	clPoolTwo := s.PrepareConcentratedPool()

	// Ensure both pools start with no incentive records
	poolOneRecords, err := clKeeper.GetAllIncentiveRecordsForPool(s.Ctx, clPoolOne.GetId())
	s.Require().NoError(err)
	s.Require().Equal(emptyIncentiveRecords, poolOneRecords)

	poolTwoRecords, err := clKeeper.GetAllIncentiveRecordsForPool(s.Ctx, clPoolTwo.GetId())
	s.Require().NoError(err)
	s.Require().Equal(emptyIncentiveRecords, poolTwoRecords)

	// Ensure setting and getting a single record works with single Get and GetAll
	clKeeper.SetIncentiveRecord(s.Ctx, incentiveRecordOne)
	poolOneRecord, err := clKeeper.GetIncentiveRecord(s.Ctx, clPoolOne.GetId(), incentiveRecordOne.IncentiveDenom, incentiveRecordOne.MinUptime)
	s.Require().NoError(err)
	s.Require().Equal(incentiveRecordOne, poolOneRecord)
	allRecordsPoolOne, err := clKeeper.GetAllIncentiveRecordsForPool(s.Ctx, clPoolOne.GetId())
	s.Require().NoError(err)
	s.Require().Equal([]types.IncentiveRecord{incentiveRecordOne}, allRecordsPoolOne)

	// Ensure records for other pool remain unchanged
	poolTwoRecord, err := clKeeper.GetIncentiveRecord(s.Ctx, clPoolTwo.GetId(), incentiveRecordOne.IncentiveDenom, incentiveRecordOne.MinUptime)
	s.Require().Error(err)
	s.Require().ErrorIs(err, types.IncentiveRecordNotFoundError{PoolId: clPoolTwo.GetId(), IncentiveDenom: incentiveRecordOne.IncentiveDenom, MinUptime: incentiveRecordOne.MinUptime})
	s.Require().Equal(types.IncentiveRecord{}, poolTwoRecord)
	allRecordsPoolTwo, err := clKeeper.GetAllIncentiveRecordsForPool(s.Ctx, clPoolTwo.GetId())
	s.Require().NoError(err)
	s.Require().Equal(emptyIncentiveRecords, allRecordsPoolTwo)

	// Ensure directly setting additional records don't overwrite previous ones
	clKeeper.SetIncentiveRecord(s.Ctx, incentiveRecordTwo)
	poolOneRecord, err = clKeeper.GetIncentiveRecord(s.Ctx, clPoolOne.GetId(), incentiveRecordTwo.IncentiveDenom, incentiveRecordTwo.MinUptime)
	s.Require().NoError(err)
	s.Require().Equal(incentiveRecordTwo, poolOneRecord)
	allRecordsPoolOne, err = clKeeper.GetAllIncentiveRecordsForPool(s.Ctx, clPoolOne.GetId())
	s.Require().NoError(err)
	s.Require().Equal([]types.IncentiveRecord{incentiveRecordOne, incentiveRecordTwo}, allRecordsPoolOne)

	// Ensure setting multiple records through helper functions as expected
	clKeeper.SetMultipleIncentiveRecords(s.Ctx, []types.IncentiveRecord{incentiveRecordThree, incentiveRecordFour})

	// Note: we expect the records to be retrieved in lexicographic order by denom
	allRecordsPoolOne, err = clKeeper.GetAllIncentiveRecordsForPool(s.Ctx, clPoolOne.GetId())
	s.Require().NoError(err)
	s.Require().Equal([]types.IncentiveRecord{incentiveRecordFour, incentiveRecordOne, incentiveRecordThree, incentiveRecordTwo}, allRecordsPoolOne)

	// Finally, we ensure the second pool remains unaffected
	allRecordsPoolTwo, err = clKeeper.GetAllIncentiveRecordsForPool(s.Ctx, clPoolTwo.GetId())
	s.Require().NoError(err)
	s.Require().Equal(emptyIncentiveRecords, allRecordsPoolTwo)
}

func (s *KeeperTestSuite) TestGetUptimeGrowthInside() {
	defaultPoolId := uint64(1)
	defaultInitialLiquidity := sdk.OneDec()
	uptimeHelper := getExpectedUptimes()

	type uptimeGrowthOutsideTest struct {
		poolSetup bool

		lowerTick                    int64
		upperTick                    int64
		currentTick                  int64
		lowerTickUptimeGrowthOutside []sdk.DecCoins
		upperTickUptimeGrowthOutside []sdk.DecCoins
		globalUptimeGrowth           []sdk.DecCoins

		expectedUptimeGrowthInside []sdk.DecCoins
		invalidTick                bool
		expectedError              bool
	}

	tests := map[string]uptimeGrowthOutsideTest{
		// current tick above range

		"current tick > upper tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  2,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is above range, we expect upper - lower
			expectedUptimeGrowthInside: uptimeHelper.hundredTokensMultiDenom,
			expectedError:              false,
		},
		"current tick > upper tick, nonzero uptime growth inside (wider range)": {
			poolSetup:                    true,
			lowerTick:                    12444,
			upperTick:                    15013,
			currentTick:                  50320,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is above range, we expect upper - lower
			expectedUptimeGrowthInside: uptimeHelper.hundredTokensMultiDenom,
			expectedError:              false,
		},
		"current tick > upper tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  2,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is above range, we expect upper - lower
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},
		"current tick > upper tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  2,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since current tick is above range, we expect upper - lower
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},
		"current tick > upper tick, zero uptime growth inside with extraneous uptime growth": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  2,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since current tick is above range, we expect upper - lower
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},

		// current tick within range

		"upper tick > current tick > lower tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    2,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is within range, we expect global - upper - lower
			expectedUptimeGrowthInside: uptimeHelper.hundredTokensMultiDenom,
			expectedError:              false,
		},
		"upper tick > current tick > lower tick, nonzero uptime growth inside (wider range)": {
			poolSetup:                    true,
			lowerTick:                    -19753,
			upperTick:                    8921,
			currentTick:                  -97,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is within range, we expect global - upper - lower
			expectedUptimeGrowthInside: uptimeHelper.hundredTokensMultiDenom,
			expectedError:              false,
		},
		"upper tick > current tick > lower tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    2,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since current tick is within range, we expect global - upper - lower
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},
		"upper tick > current tick > lower tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    2,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since current tick is within range, we expect global - upper - lower
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},

		// current tick below range

		"current tick < lower tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  -1,
			lowerTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is below range, we expect lower - upper
			expectedUptimeGrowthInside: uptimeHelper.hundredTokensMultiDenom,
			expectedError:              false,
		},
		"current tick < lower tick, nonzero uptime growth inside (wider range)": {
			poolSetup:                    true,
			lowerTick:                    328,
			upperTick:                    726,
			currentTick:                  189,
			lowerTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is below range, we expect lower - upper
			expectedUptimeGrowthInside: uptimeHelper.hundredTokensMultiDenom,
			expectedError:              false,
		},
		"current tick < lower tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  -1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is below range, we expect lower - upper
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},
		"current tick < lower tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  -1,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since current tick is below range, we expect lower - upper
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},

		// current tick on range boundary

		"current tick = lower tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  0,
			lowerTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.fourHundredTokensMultiDenom,

			// Since we treat the range as [lower, upper) (i.e. inclusive of lower tick, exclusive of upper),
			// this case is equivalent to the current tick being within the range (global - upper - lower)
			expectedUptimeGrowthInside: uptimeHelper.hundredTokensMultiDenom,
			expectedError:              false,
		},
		"current tick = lower tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  0,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since we treat the range as [lower, upper) (i.e. inclusive of lower tick, exclusive of upper),
			// this case is equivalent to the current tick being within the range (global - upper - lower)
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},
		"current tick = lower tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  0,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since we treat the range as [lower, upper) (i.e. inclusive of lower tick, exclusive of upper),
			// this case is equivalent to the current tick being within the range (global - upper - lower)
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},
		"current tick = upper tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.fourHundredTokensMultiDenom,

			// Since we treat the range as [lower, upper) (i.e. inclusive of lower tick, exclusive of upper),
			// this case is equivalent to the current tick being above the range (upper - lower)
			expectedUptimeGrowthInside: uptimeHelper.hundredTokensMultiDenom,
			expectedError:              false,
		},
		"current tick = upper tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since we treat the range as [lower, upper) (i.e. inclusive of lower tick, exclusive of upper),
			// this case is equivalent to the current tick being above the range (upper - lower)
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},
		"current tick = upper tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since we treat the range as [lower, upper) (i.e. inclusive of lower tick, exclusive of upper),
			// this case is equivalent to the current tick being above the range (upper - lower)
			expectedUptimeGrowthInside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:              false,
		},

		// error catching

		"error: pool has not been setup": {
			poolSetup:     false,
			expectedError: true,
		},
	}

	for name, tc := range tests {
		s.Run(name, func() {
			s.SetupTest()

			// if pool set up true, set up default pool
			var pool types.ConcentratedPoolExtension
			if tc.poolSetup {
				pool = s.PrepareConcentratedPool()
				currentTick := pool.GetCurrentTick().Int64()

				// Update global uptime accums
				addToUptimeAccums(s.Ctx, pool.GetId(), s.App.ConcentratedLiquidityKeeper, tc.globalUptimeGrowth)

				// Update tick-level uptime trackers
				s.initializeTick(s.Ctx, currentTick, tc.lowerTick, defaultInitialLiquidity, cl.EmptyCoins, wrapUptimeTrackers(tc.lowerTickUptimeGrowthOutside), true)
				s.initializeTick(s.Ctx, currentTick, tc.upperTick, defaultInitialLiquidity, cl.EmptyCoins, wrapUptimeTrackers(tc.upperTickUptimeGrowthOutside), false)
				pool.SetCurrentTick(sdk.NewInt(tc.currentTick))
				s.App.ConcentratedLiquidityKeeper.SetPool(s.Ctx, pool)
			}

			// system under test
			uptimeGrowthInside, err := s.App.ConcentratedLiquidityKeeper.GetUptimeGrowthInsideRange(s.Ctx, defaultPoolId, tc.lowerTick, tc.upperTick)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				// check if returned uptime growth inside has correct value
				s.Require().Equal(tc.expectedUptimeGrowthInside, uptimeGrowthInside)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGetUptimeGrowthOutside() {
	defaultPoolId := uint64(1)
	defaultInitialLiquidity := sdk.OneDec()
	uptimeHelper := getExpectedUptimes()

	type uptimeGrowthOutsideTest struct {
		poolSetup bool

		lowerTick                    int64
		upperTick                    int64
		currentTick                  int64
		lowerTickUptimeGrowthOutside []sdk.DecCoins
		upperTickUptimeGrowthOutside []sdk.DecCoins
		globalUptimeGrowth           []sdk.DecCoins

		expectedUptimeGrowthOutside []sdk.DecCoins
		invalidTick                 bool
		expectedError               bool
	}

	tests := map[string]uptimeGrowthOutsideTest{
		// current tick above range

		"current tick > upper tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  2,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is above range, we expect global - (upper - lower)
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick > upper tick, nonzero uptime growth inside (wider range)": {
			poolSetup:                    true,
			lowerTick:                    12444,
			upperTick:                    15013,
			currentTick:                  50320,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is above range, we expect global - (upper - lower)
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick > upper tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  2,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.threeHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick > upper tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  2,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:               false,
		},
		"current tick > upper tick, zero uptime growth inside with extraneous uptime growth": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  2,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},

		// current tick within range

		"upper tick > current tick > lower tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    2,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is within range, we expect global - (global - upper - lower)
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"upper tick > current tick > lower tick, nonzero uptime growth inside (wider range)": {
			poolSetup:                    true,
			lowerTick:                    -19753,
			upperTick:                    8921,
			currentTick:                  -97,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since current tick is within range, we expect global - (global - upper - lower)
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"upper tick > current tick > lower tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    2,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"upper tick > current tick > lower tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    2,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:               false,
		},

		// current tick below range

		"current tick < lower tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  -1,
			lowerTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is below range, we expect global - (lower - upper)
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick < lower tick, nonzero uptime growth inside (wider range)": {
			poolSetup:                    true,
			lowerTick:                    328,
			upperTick:                    726,
			currentTick:                  189,
			lowerTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since current tick is below range, we expect global - (lower - upper)
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick < lower tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  -1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.threeHundredTokensMultiDenom,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.threeHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick < lower tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  -1,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:               false,
		},

		// current tick on range boundary

		"current tick = lower tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  0,
			lowerTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.fourHundredTokensMultiDenom,

			// Since we treat the range as [lower, upper) (i.e. inclusive of lower tick, exclusive of upper),
			// this case is equivalent to the current tick being within the range (global - (global - upper - lower))
			expectedUptimeGrowthOutside: uptimeHelper.threeHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick = lower tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  0,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick = lower tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  0,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:               false,
		},
		"current tick = upper tick, nonzero uptime growth inside": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.fourHundredTokensMultiDenom,

			// Since we treat the range as [lower, upper) (i.e. inclusive of lower tick, exclusive of upper),
			// this case is equivalent to the current tick being above the range (global - (upper - lower))
			expectedUptimeGrowthOutside: uptimeHelper.threeHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick = upper tick, zero uptime growth inside (nonempty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			upperTickUptimeGrowthOutside: uptimeHelper.hundredTokensMultiDenom,
			globalUptimeGrowth:           uptimeHelper.twoHundredTokensMultiDenom,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.twoHundredTokensMultiDenom,
			expectedError:               false,
		},
		"current tick = upper tick, zero uptime growth inside (empty trackers)": {
			poolSetup:                    true,
			lowerTick:                    0,
			upperTick:                    1,
			currentTick:                  1,
			lowerTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			upperTickUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			globalUptimeGrowth:           uptimeHelper.emptyExpectedAccumValues,

			// Since the range is empty, we expect growth outside to be equal to global
			expectedUptimeGrowthOutside: uptimeHelper.emptyExpectedAccumValues,
			expectedError:               false,
		},

		// error catching

		"error: pool has not been setup": {
			poolSetup:     false,
			expectedError: true,
		},
	}

	for name, tc := range tests {
		s.Run(name, func() {
			s.SetupTest()

			// if pool set up true, set up default pool
			var pool types.ConcentratedPoolExtension
			if tc.poolSetup {
				pool = s.PrepareConcentratedPool()
				currentTick := pool.GetCurrentTick().Int64()

				// Update global uptime accums
				addToUptimeAccums(s.Ctx, pool.GetId(), s.App.ConcentratedLiquidityKeeper, tc.globalUptimeGrowth)

				// Update tick-level uptime trackers
				s.initializeTick(s.Ctx, currentTick, tc.lowerTick, defaultInitialLiquidity, cl.EmptyCoins, wrapUptimeTrackers(tc.lowerTickUptimeGrowthOutside), true)
				s.initializeTick(s.Ctx, currentTick, tc.upperTick, defaultInitialLiquidity, cl.EmptyCoins, wrapUptimeTrackers(tc.upperTickUptimeGrowthOutside), false)
				pool.SetCurrentTick(sdk.NewInt(tc.currentTick))
				s.App.ConcentratedLiquidityKeeper.SetPool(s.Ctx, pool)
			}

			// system under test
			uptimeGrowthOutside, err := s.App.ConcentratedLiquidityKeeper.GetUptimeGrowthOutsideRange(s.Ctx, defaultPoolId, tc.lowerTick, tc.upperTick)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				// check if returned uptime growth inside has correct value
				s.Require().Equal(tc.expectedUptimeGrowthOutside, uptimeGrowthOutside)
			}
		})
	}
}