<?php

namespace Tests\Unit\Players;

use App\Models\DatabaseHandler;
use App\Players\PlayerDoctrineRepository;
use Tests\Helper\ModelHelper;
use Tests\Helper\PlayerHelper;
use Tests\TestCase;

/**
 * Class PlayerDoctrineRepositoryTest
 *
 * @package Tests\Unit\Players
 */
final class PlayerDoctrineRepositoryTest extends TestCase
{
    use ModelHelper;
    use PlayerHelper;

    //region Tests

    /**
     * @param bool $withPlayer
     *
     * @return array
     */
    private function setUpFindOneByNameTest(bool $withPlayer = true): array
    {
        $name = $this->getFaker()->word;
        $player = $this->createPlayerModel();
        $databaseHandler = $this->createDatabaseHandler();
        $this->mockDatabaseHandlerLoad($databaseHandler, $withPlayer ? $player : null, ['name' => $name]);
        $playerRepository = $this->getPlayerDoctrineRepository($databaseHandler);

        return [$playerRepository, $name, $player];
    }

    /**
     * @return void
     */
    public function testFindOneByName(): void
    {
        /** @var PlayerDoctrineRepository $playerRepository */
        [$playerRepository, $name, $player] = $this->setUpFindOneByNameTest();

        $this->assertEquals($player, $playerRepository->findOneByName($name));
    }

    /**
     * @return void
     */
    public function testFindOneByNameWithoutPlayer(): void
    {
        /** @var PlayerDoctrineRepository $playerRepository */
        [$playerRepository, $name] = $this->setUpFindOneByNameTest(false);

        $this->assertNull($playerRepository->findOneByName($name));
    }

    //endregion

    /**
     * @param DatabaseHandler|null $databaseHandler
     *
     * @return PlayerDoctrineRepository
     */
    private function getPlayerDoctrineRepository(DatabaseHandler $databaseHandler = null): PlayerDoctrineRepository
    {
        return new PlayerDoctrineRepository($databaseHandler ?: $this->createDatabaseHandler());
    }
}
