<?php

namespace Tests\Unit\Seasons;

use App\Models\Exceptions\ModelNotFoundException;
use App\Seasons\SeasonManager;
use App\Seasons\SeasonModel;
use App\Seasons\SeasonModelFactory;
use App\Seasons\SeasonRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Tests\Helper\ModelHelper;
use Tests\Helper\SeasonHelper;
use Tests\TestCase;

/**
 * Class SeasonManagerTest
 *
 * @package Tests\Unit\Seasons
 */
final class SeasonManagerTest extends TestCase
{
    use ModelHelper;
    use SeasonHelper;

    //region Tests

    /**
     * @return void
     */
    public function testCreateSeason(): void
    {
        $name = $this->getFaker()->word;
        $season = $this->createSeasonModel();
        $seasonModelFactory = $this->createSeasonModelFactory();
        $this->mockSeasonModelFactoryCreate($seasonModelFactory, $season, $name);
        $seasonRepository = $this->createSeasonRepository();
        $this->mockRepositorySave($seasonRepository, $season);

        $this->assertEquals($season, $this->getSeasonManager($seasonRepository, $seasonModelFactory)->createSeason($name));
        $this->assertRepositorySave($seasonRepository, $season);
    }

    /**
     * @param bool $withSeason
     *
     * @return array
     */
    private function setUpGetSeasonByNameTest(bool $withSeason = true): array
    {
        $name = $this->getFaker()->word;
        $season = $this->createSeasonModel();
        $seasonRepository = $this->createSeasonRepository();
        $this->mockSeasonRepositoryFindOneByName($seasonRepository, $withSeason ? $season : null, $name);
        $seasonManager = $this->getSeasonManager($seasonRepository);

        return [$seasonManager, $name, $season];
    }

    /**
     * @return void
     */
    public function testGetSeasonByName(): void
    {
        /** @var SeasonManager $seasonManager */
        [$seasonManager, $name, $season] = $this->setUpGetSeasonByNameTest();

        $this->assertEquals($season, $seasonManager->getSeasonByName($name));
    }

    /**
     * @return void
     */
    public function testGetSeasonByNameWithoutSeason(): void
    {
        /** @var SeasonManager $seasonManager */
        [$seasonManager, $name] = $this->setUpGetSeasonByNameTest(false);

        $this->expectException(ModelNotFoundException::class);

        $seasonManager->getSeasonByName($name);
    }

    /**
     * @return void
     */
    public function testGetSeasons(): void
    {
        $season = $this->createSeasonModel();
        $seasonRepository = $this->createSeasonRepository();
        $this->mockRepositoryFindAll($seasonRepository, new ArrayCollection([$season]));

        $this->assertEquals(
            new ArrayCollection([$season]),
            $this->getSeasonManager($seasonRepository)->getSeasons()
        );
    }

    /**
     * @return array
     */
    private function setUpActivateSeasonTest(): array
    {
        $season = $this->createSeasonModel();
        $season
            ->shouldReceive('setActive')
            ->with(true)
            ->andReturn($season)
            ->once();
        $seasonRepository = $this->createSeasonRepository();
        $seasonRepository
            ->shouldReceive('deactivateActiveSeason')
            ->andReturn($seasonRepository)
            ->once();
        $this->mockRepositorySave($seasonRepository, $season);
        $seasonManager = $this->getSeasonManager($seasonRepository);

        return [$seasonManager, $season];
    }

    /**
     * @return void
     */
    public function testActivateSeason(): void
    {
        /**
         * @var SeasonManager $seasonManager
         * @var SeasonModel   $season
         */
        [$seasonManager, $season] = $this->setUpActivateSeasonTest();

        $this->assertEquals($season, $seasonManager->activateSeason($season));
    }

    //endregion

    private function getSeasonManager(
        SeasonRepository $seasonRepository = null,
        SeasonModelFactory $seasonModelFactory = null
    ): SeasonManager {
        return new SeasonManager(
            $seasonRepository ?: $this->createSeasonRepository(),
            $seasonModelFactory ?: $this->createSeasonModelFactory()
        );
    }
}
