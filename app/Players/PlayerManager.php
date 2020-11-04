<?php

namespace App\Players;

use App\Models\Exceptions\ModelNotFoundException;

/**
 * Class PlayerManager
 *
 * @package App\Players
 */
class PlayerManager
{
    /**
     * @var PlayerRepository
     */
    private PlayerRepository $playerRepository;

    /**
     * @var PlayerModelFactory
     */
    private PlayerModelFactory $playerModelFactory;

    /**
     * PlayerManager constructor.
     *
     * @param PlayerRepository   $playerRepository
     * @param PlayerModelFactory $playerModelFactory
     */
    public function __construct(
        PlayerRepository $playerRepository,
        PlayerModelFactory $playerModelFactory
    ) {
        $this->playerRepository = $playerRepository;
        $this->playerModelFactory = $playerModelFactory;
    }

    /**
     * @return PlayerRepository
     */
    private function getPlayerRepository(): PlayerRepository
    {
        return $this->playerRepository;
    }

    /**
     * @return PlayerModelFactory
     */
    private function getPlayerModelFactory(): PlayerModelFactory
    {
        return $this->playerModelFactory;
    }

    /**
     * @param string $name
     *
     * @return PlayerModel
     */
    public function createPlayer(string $name): PlayerModel
    {
        $player = $this->getPlayerModelFactory()->create($name);

        return $this->getPlayerRepository()->save($player);
    }

    /**
     * @param string $name
     *
     * @return PlayerModel
     */
    public function getPlayerByName(string $name): PlayerModel
    {
        $player = $this->getPlayerRepository()->findOneByName($name);
        if (!$player) {
            throw new ModelNotFoundException(\sprintf('Player with name %s not found', $name));
        }

        return $player;
    }
}
