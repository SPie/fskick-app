<?php

namespace App\Players;

/**
 * Class PlayerDoctrineModelFactory
 *
 * @package App\Players
 */
final class PlayerDoctrineModelFactory implements PlayerModelFactory
{
    /**
     * @param string $name
     *
     * @return PlayerModel
     */
    public function create(string $name): PlayerModel
    {
        return new PlayerDoctrineModel($name);
    }
}
