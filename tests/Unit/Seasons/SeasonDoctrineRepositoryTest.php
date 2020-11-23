<?php

namespace Tests\Unit\Seasons;

use App\Models\DatabaseHandler;
use App\Models\Exceptions\ModelNotFoundException;
use App\Seasons\SeasonDoctrineRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Tests\Helper\ModelHelper;
use Tests\Helper\SeasonHelper;
use Tests\TestCase;

/**
 * Class SeasonDoctrineRepositoryTest
 *
 * @package Tests\Unit\Seasons
 */
final class SeasonDoctrineRepositoryTest extends TestCase
{
    use ModelHelper;
    use SeasonHelper;

    //region Tests

    /**
     * @param bool $withSeason
     *
     * @return array
     */
    private function setUpFindOneByNameTest(bool $withSeason = true): array
    {
        $name = $this->getFaker()->word;
        $season = $this->createSeasonModel();
        $databaseHandler = $this->createDatabaseHandler();
        $this->mockDatabaseHandlerLoad($databaseHandler, $withSeason ? $season : null, ['name' => $name]);
        $seasonRepository = $this->getSeasonDoctrineRepository($databaseHandler);

        return [$seasonRepository, $name, $season];
    }

    /**
     * @return void
     */
    public function testFindOneByName(): void
    {
        /** @var SeasonDoctrineRepository $seasonRepository */
        [$seasonRepository, $name, $season] = $this->setUpFindOneByNameTest();

        $this->assertEquals($season, $seasonRepository->findOneByName($name));
    }

    /**
     * @return void
     */
    public function testFindOneByNameWithoutSeason(): void
    {
        /** @var SeasonDoctrineRepository $seasonRepository */
        [$seasonRepository, $name] = $this->setUpFindOneByNameTest(false);

        $this->assertNull($seasonRepository->findOneByName($name));
    }

    /**
     * @return void
     */
    public function testDeactivateActiveSeason(): void
    {
        $season1 = $this->createSeasonModel();
        $season1
            ->shouldReceive('setActive')
            ->with(true)
            ->andReturn($season1)
            ->once();
        $season2 = $this->createSeasonModel();
        $season2
            ->shouldReceive('setActive')
            ->with(true)
            ->andReturn($season1)
            ->once();
        $databaseHandler = $this->createDatabaseHandler();
        $databaseHandler
            ->shouldReceive('loadAll')
            ->with(['active' => true], [], null, null)
            ->andReturn(new ArrayCollection([$season1, $season2]));
        $this->mockDatabaseHandlerSave($databaseHandler, $season1, false);
        $this->mockDatabaseHandlerSave($databaseHandler, $season2, false);
        $seasonRepository = $this->getSeasonDoctrineRepository($databaseHandler);

        $seasonRepository->deactivateActiveSeason();

        $this->assertDatabaseHandlerFlush($databaseHandler);
    }

    //endregion

    /**
     * @param DatabaseHandler|null $databaseHandler
     *
     * @return SeasonDoctrineRepository
     */
    private function getSeasonDoctrineRepository(DatabaseHandler $databaseHandler = null): SeasonDoctrineRepository
    {
        return new SeasonDoctrineRepository($databaseHandler ?: $this->createDatabaseHandler());
    }
}
