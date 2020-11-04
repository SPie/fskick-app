<?php

namespace App\Seasons;

/**
 * Class SeasonDoctrineModelFactory
 *
 * @package App\Seasons
 */
final class SeasonDoctrineModelFactory implements SeasonModelFactory
{
    /**
     * @param string $name
     *
     * @return SeasonModel
     */
    public function create(string $name): SeasonModel
    {
        return new SeasonDoctrineModel($name);
    }
}
