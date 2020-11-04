<?php

namespace Tests\Unit\Seasons;

use App\Models\DatabaseHandler;
use App\Models\Exceptions\ModelNotFoundException;
use App\Seasons\SeasonDoctrineRepository;
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
