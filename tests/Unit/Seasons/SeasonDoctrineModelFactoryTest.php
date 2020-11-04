<?php

namespace Tests\Unit\Seasons;

use App\Seasons\SeasonDoctrineModel;
use App\Seasons\SeasonDoctrineModelFactory;
use Tests\TestCase;

/**
 * Class SeasonDoctrineModelFactoryTest
 *
 * @package Tests\Unit\Seasons
 */
final class SeasonDoctrineModelFactoryTest extends TestCase
{
    //region Tests

    /**
     * @return void
     */
    public function testCreate(): void
    {
        $name = $this->getFaker()->word;

        $this->assertEquals(
            new SeasonDoctrineModel($name),
            $this->getSeasonDoctrineModelFactory()->create($name)
        );
    }

    //endregion

    /**
     * @return SeasonDoctrineModelFactory
     */
    private function getSeasonDoctrineModelFactory(): SeasonDoctrineModelFactory
    {
        return new SeasonDoctrineModelFactory();
    }
}
