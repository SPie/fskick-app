<?php

namespace Tests\Helper;

use App\Players\PlayerDoctrineModel;
use App\Players\PlayerManager;
use App\Players\PlayerModel;
use App\Players\PlayerModelFactory;
use App\Players\PlayerRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Mockery as m;
use Mockery\MockInterface;

/**
 * Trait PlayerHelper
 *
 * @package Tests\Helper
 */
trait PlayerHelper
{
    /**
     * @return PlayerModel|MockInterface
     */
    private function createPlayerModel(): PlayerModel
    {
        return m::spy(PlayerModel::class);
    }

    /**
     * @param PlayerModel|MockInterface $playerModel
     * @param string                    $name
     *
     * @return $this
     */
    private function mockPlayerModelGetName(MockInterface $playerModel, string $name): self
    {
        $playerModel
            ->shouldReceive('getName')
            ->andReturn($name);

        return $this;
    }

    /**
     * @param int   $times
     * @param array $attributes
     *
     * @return PlayerModel[]|Collection
     */
    private function createPlayerEntities(int $times = 1, array $attributes = []): Collection
    {
        if ($times = 1) {
            return new ArrayCollection([entity(PlayerDoctrineModel::class, 1)->create($attributes)]);
        }

        return entity(PlayerDoctrineModel::class, $times)->create($attributes);
    }

    /**
     * @return PlayerModelFactory|MockInterface
     */
    private function createPlayerModelFactory(): PlayerModelFactory
    {
        return m::spy(PlayerModelFactory::class);
    }

    /**
     * @param PlayerModelFactory|MockInterface $playerModelFactory
     * @param PlayerModel                      $player
     * @param string                           $name
     *
     * @return $this
     */
    private function mockPlayerModelFactoryCreate(MockInterface $playerModelFactory, PlayerModel $player, string $name): self
    {
        $playerModelFactory
            ->shouldReceive('create')
            ->with($name)
            ->andReturn($player);

        return $this;
    }

    /**
     * @return PlayerRepository|MockInterface
     */
    private function createPlayerRepository(): PlayerRepository
    {
        return m::spy(PlayerRepository::class);
    }

    /**
     * @param PlayerRepository|MockInterface $playerRepository
     * @param PlayerModel|null               $player
     * @param string                         $name
     *
     * @return $this
     */
    private function mockPlayerRepositoryFindOneByName(
        MockInterface $playerRepository,
        ?PlayerModel $player,
        string $name
    ): self {
        $playerRepository
            ->shouldReceive('findOneByName')
            ->with($name)
            ->andReturn($player);

        return $this;
    }

    /**
     * @return PlayerManager
     */
    private function createPlayerManager(): PlayerManager
    {
        return m::spy(PlayerManager::class);
    }

    /**
     * @param PlayerManager|MockInterface $playerManager
     * @param PlayerModel                 $player
     * @param string                      $name
     *
     * @return $this
     */
    private function mockPlayerManagerCreatePlayer(MockInterface $playerManager, PlayerModel $player, string $name): self
    {
        $playerManager
            ->shouldReceive('createPlayer')
            ->with($name)
            ->andReturn($player);

        return $this;
    }

    /**
     * @param PlayerManager|MockInterface $playerManager
     * @param string                      $name
     *
     * @return $this
     */
    private function assertPlayerManagerCreatePlayer(MockInterface $playerManager, string $name): self
    {
        $playerManager
            ->shouldHaveReceived('createPlayer')
            ->with($name)
            ->once();

        return $this;
    }

    /**
     * @param PlayerManager|MockInterface $playerManager
     * @param PlayerModel|\Exception      $player
     * @param string                      $name
     *
     * @return $this
     */
    private function mockPlayerManagerGetPlayerByName(MockInterface $playerManager, $player, string $name): self
    {
        $playerManager
            ->shouldReceive('getPlayerByName')
            ->with($name)
            ->andThrow($player);

        return $this;
    }
}
