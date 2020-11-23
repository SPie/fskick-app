<?php

namespace App\Seasons;

use App\Models\AbstractDoctrineRepository;
use App\Models\Model;

/**
 * Class SeasonDoctrineRepository
 *
 * @package App\Seasons
 */
final class SeasonDoctrineRepository extends AbstractDoctrineRepository implements SeasonRepository
{
    /**
     * @param string $name
     *
     * @return SeasonModel|Model|null
     */
    public function findOneByName(string $name): ?SeasonModel
    {
        return $this->findOneBy([SeasonModel::PROPERTY_NAME => $name]);
    }

    /**
     * @return SeasonRepository
     */
    public function deactivateActiveSeason(): SeasonRepository
    {
        $activeSeasons = $this->findBy([SeasonModel::PROPERTY_ACTIVE => true]);
        foreach ($activeSeasons as $activeSeason) {
            $this->save($activeSeason->setActive(false), false);
        }

        $this->flush();

        return $this;
    }
}
