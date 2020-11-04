<?php

namespace App\Seasons;

/**
 * Interface SeasonModelFactory
 *
 * @package App\Seasons
 */
interface SeasonModelFactory
{
    /**
     * @param string $name
     *
     * @return SeasonModel
     */
    public function create(string $name): SeasonModel;
}
