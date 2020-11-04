<?php

namespace App\Players;

use App\Models\Model;
use App\Models\Repository;

/**
 * Interface PlayerRepository
 *
 * @package App\Players
 */
interface PlayerRepository extends Repository
{
    /**
     * @param Model $model
     * @param bool  $flush
     *
     * @return PlayerModel|Model
     */
    public function save(Model $model, bool $flush = true): Model;

    /**
     * @param string $name
     *
     * @return PlayerModel|null
     */
    public function findOneByName(string $name): ?PlayerModel;
}
