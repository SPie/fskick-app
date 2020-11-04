<?php

namespace Tests\Unit\Players;

use App\Players\PlayerDoctrineModel;
use App\Players\PlayerDoctrineModelFactory;
use Tests\TestCase;

/**
 * Class PlayerDoctrineModelFactoryTest
 *
 * @package Tests\Unit\Players
 */
final class PlayerDoctrineModelFactoryTest extends TestCase
{
    //region Tests

    /**
     * @return void
     */
    public function testCreate(): void
    {
        $name = $this->getFaker()->word;

        $this->assertEquals(
            new PlayerDoctrineModel($name),
            $this->getPlayerDoctrineModelFactory()->create($name)
        );
    }

    //endregion

    /**
     * @return PlayerDoctrineModelFactory
     */
    private function getPlayerDoctrineModelFactory(): PlayerDoctrineModelFactory
    {
        return new PlayerDoctrineModelFactory();
    }
}
