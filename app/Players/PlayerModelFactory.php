<?php

namespace App\Players;

/**
 * Interface PlayerModelFactory
 *
 * @package App\Players
 */
interface PlayerModelFactory
{
    /**
     * @param string $name
     *
     * @return PlayerModel
     */
    public function create(string $name): PlayerModel;
}
