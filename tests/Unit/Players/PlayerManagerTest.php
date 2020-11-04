<?php

namespace Tests\Unit\Players;

use App\Models\Exceptions\ModelNotFoundException;
use App\Players\PlayerManager;
use App\Players\PlayerModelFactory;
use App\Players\PlayerRepository;
use Tests\Helper\ModelHelper;
use Tests\Helper\PlayerHelper;
use Tests\TestCase;

/**
 * Class PlayerManagerTest
 *
 * @package Tests\Unit
 */
final class PlayerManagerTest extends TestCase
{
    use ModelHelper;
    use PlayerHelper;

    //region Tests

    /**
     * @return void
     */
    public function testCreatePlayer(): void
    {
        $name = $this->getFaker()->word;
        $player = $this->createPlayerModel();
        $playerModelFactory = $this->createPlayerModelFactory();
        $this->mockPlayerModelFactoryCreate($playerModelFactory, $player, $name);
        $playerRepository = $this->createPlayerRepository();
        $this->mockRepositorySave($playerRepository, $player);

        $this->assertEquals($player, $this->getPlayerManager($playerRepository, $playerModelFactory)->createPlayer($name));
        $this->assertRepositorySave($playerRepository, $player);
    }

    /**
     * @param bool $withPlayer
     *
     * @return array
     */
    private function setUpGetPlayerByNameTest(bool $withPlayer = true): array
    {
        $name = $this->getFaker()->word;
        $player = $this->createPlayerModel();
        $playerRepository = $this->createPlayerRepository();
        $this->mockPlayerRepositoryFindOneByName($playerRepository, $withPlayer ? $player : null, $name);
        $playerManager = $this->getPlayerManager($playerRepository);

        return [$playerManager, $name, $player];
    }

    /**
     * @return void
     */
    public function testGetPlayerByName(): void
    {
        /** @var PlayerManager $playerManager */
        [$playerManager, $name, $player] = $this->setUpGetPlayerByNameTest();

        $this->assertEquals($player, $playerManager->getPlayerByName($name));
    }

    /**
     * @return void
     */
    public function testGetPlayerByNameWithoutPlayer(): void
    {
        /** @var PlayerManager $playerManager */
        [$playerManager, $name] = $this->setUpGetPlayerByNameTest(false);

        $this->expectException(ModelNotFoundException::class);

        $playerManager->getPlayerByName($name);
    }

    //endregion

    /**
     * @param PlayerRepository|null   $playerRepository
     * @param PlayerModelFactory|null $playerModelFactory
     *
     * @return PlayerManager
     */
    private function getPlayerManager(
        PlayerRepository $playerRepository = null,
        PlayerModelFactory $playerModelFactory = null
    ): PlayerManager {
        return new PlayerManager(
            $playerRepository ?: $this->createPlayerRepository(),
            $playerModelFactory ?: $this->createPlayerModelFactory()
        );
    }
}
