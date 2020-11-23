<?php

namespace App\Seasons;

use App\Models\Exceptions\ModelNotFoundException;
use Doctrine\Common\Collections\Collection;

/**
 * Class SeasonManager
 *
 * @package App\Seasons
 */
class SeasonManager
{
    /**
     * @var SeasonRepository
     */
    private SeasonRepository $seasonRepository;

    /**
     * @var SeasonModelFactory
     */
    private SeasonModelFactory $seasonModelFactory;

    /**
     * SeasonManager constructor.
     *
     * @param SeasonRepository   $seasonRepository
     * @param SeasonModelFactory $seasonModelFactory
     */
    public function __construct(SeasonRepository $seasonRepository, SeasonModelFactory $seasonModelFactory)
    {
        $this->seasonRepository = $seasonRepository;
        $this->seasonModelFactory = $seasonModelFactory;
    }

    /**
     * @return SeasonRepository
     */
    private function getSeasonRepository(): SeasonRepository
    {
        return $this->seasonRepository;
    }

    /**
     * @return SeasonModelFactory
     */
    private function getSeasonModelFactory(): SeasonModelFactory
    {
        return $this->seasonModelFactory;
    }

    /**
     * @param string $name
     *
     * @return SeasonModel
     */
    public function createSeason(string $name): SeasonModel
    {
        $season = $this->getSeasonModelFactory()->create($name);

        return $this->getSeasonRepository()->save($season);
    }

    /**
     * @return SeasonModel[]|Collection
     */
    public function getSeasons(): Collection
    {
        return $this->getSeasonRepository()->findAll();
    }

    /**
     * @param string $name
     *
     * @return SeasonModel
     */
    public function getSeasonByName(string $name): SeasonModel
    {
        $season = $this->getSeasonRepository()->findOneByName($name);
        if (!$season) {
            throw new ModelNotFoundException(\sprintf('Season with name %s not found', $name));
        }

        return $season;
    }

    /**
     * @param SeasonModel $season
     *
     * @return SeasonModel
     */
    public function activateSeason(SeasonModel $season): SeasonModel
    {
        return $this->getSeasonRepository()
            ->deactivateActiveSeason()
            ->save($season->setActive(true));
    }
}
