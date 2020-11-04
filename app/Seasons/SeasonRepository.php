<?php

namespace App\Seasons;

use App\Models\Model;
use App\Models\Repository;

/**
 * Interface SeasonRepository
 *
 * @package App\Seasons
 */
interface SeasonRepository extends Repository
{
    /**
     * @param SeasonModel|Model $model
     * @param bool              $flush
     *
     * @return SeasonModel|Model
     */
    public function save(Model $model, bool $flush = true): Model;

    /**
     * @param string $name
     *
     * @return SeasonModel|null
     */
    public function findOneByName(string $name): ?SeasonModel;
}
