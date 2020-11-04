<?php

namespace App\Players;

use App\Models\AbstractDoctrineRepository;
use App\Models\Model;

/**
 * Class PlayerDoctrineRepository
 *
 * @package App\Players
 */
final class PlayerDoctrineRepository extends AbstractDoctrineRepository implements PlayerRepository
{
    /**
     * @param string $name
     *
     * @return PlayerModel|Model|null
     */
    public function findOneByName(string $name): ?PlayerModel
    {
        return $this->findOneBy([PlayerModel::PROPERTY_NAME => $name]);
    }
}
