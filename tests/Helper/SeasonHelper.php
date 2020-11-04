<?php

namespace Tests\Helper;

use App\Seasons\SeasonDoctrineModel;
use App\Seasons\SeasonManager;
use App\Seasons\SeasonModel;
use App\Seasons\SeasonModelFactory;
use App\Seasons\SeasonRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Mockery as m;
use Mockery\MockInterface;

/**
 * Trait SeasonHelper
 *
 * @package Tests\Helper
 */
trait SeasonHelper
{
    /**
     * @return SeasonModel|MockInterface
     */
    private function createSeasonModel(): SeasonModel
    {
        return m::spy(SeasonModel::class);
    }

    /**
     * @param SeasonModel|MockInterface $seasonModel
     * @param string                    $name
     *
     * @return $this
     */
    private function mockSeasonModelGetName(MockInterface $seasonModel, string $name): self
    {
        $seasonModel
            ->shouldReceive('getName')
            ->andReturn($name);

        return $this;
    }

    private function createSeasonEntities(int $times = 1, array $attributes = []): Collection
    {
        if ($times = 1) {
            return new ArrayCollection([entity(SeasonDoctrineModel::class, 1)->create($attributes)]);
        }

        return entity(SeasonDoctrineModel::class, $times)->create($attributes);
    }

    /**
     * @return SeasonManager|MockInterface
     */
    private function createSeasonManager(): SeasonManager
    {
        return m::spy(SeasonManager::class);
    }

    /**
     * @param SeasonManager|MockInterface $seasonManager
     * @param SeasonModel                 $season
     * @param string                      $name
     *
     * @return $this
     */
    private function mockSeasonManagerCreateSeason(MockInterface $seasonManager, SeasonModel $season, string $name): self
    {
        $seasonManager
            ->shouldReceive('createSeason')
            ->with($name)
            ->andReturn($season);

        return $this;
    }

    /**
     * @param SeasonManager|MockInterface $seasonManager
     * @param string                      $name
     *
     * @return $this
     */
    private function assertSeasonManagerCreateSeason(MockInterface $seasonManager, string $name): self
    {
        $seasonManager
            ->shouldHaveReceived('createSeason')
            ->with($name)
            ->once();

        return $this;
    }

    /**
     * @param SeasonManager|MockInterface $seasonManager
     * @param SeasonModel|\Exception      $season
     * @param string                      $name
     *
     * @return $this
     */
    private function mockSeasonManagerGetSeasonByName(MockInterface $seasonManager, $season, string $name): self
    {
        $seasonManager
            ->shouldReceive('getSeasonByName')
            ->with($name)
            ->andThrow($season);

        return $this;
    }

    /**
     * @return SeasonRepository|MockInterface
     */
    private function createSeasonRepository(): SeasonRepository
    {
        return m::spy(SeasonRepository::class);
    }

    /**
     * @param SeasonRepository|MockInterface $seasonRepository
     * @param SeasonModel|null               $season
     * @param string                         $name
     *
     * @return $this
     */
    private function mockSeasonRepositoryFindOneByName(
        MockInterface $seasonRepository,
        ?SeasonModel $season,
        string $name
    ): self {
        $seasonRepository
            ->shouldReceive('findOneByName')
            ->with($name)
            ->andReturn($season);

        return $this;
    }

    /**
     * @return SeasonModelFactory|MockInterface
     */
    private function createSeasonModelFactory(): SeasonModelFactory
    {
        return m::spy(SeasonModelFactory::class);
    }

    /**
     * @param SeasonModelFactory|MockInterface $seasonModelFactory
     * @param SeasonModel                      $season
     * @param string                           $name
     *
     * @return $this
     */
    private function mockSeasonModelFactoryCreate(MockInterface $seasonModelFactory, SeasonModel $season, string $name): self
    {
        $seasonModelFactory
            ->shouldReceive('create')
            ->with($name)
            ->andReturn($season);

        return $this;
    }
}
