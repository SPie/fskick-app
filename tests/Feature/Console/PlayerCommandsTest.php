<?php

namespace Tests\Feature\Console;

use App\Players\PlayerModel;
use App\Players\PlayerRepository;
use LaravelDoctrine\Migrations\Testing\DatabaseMigrations;
use Tests\FeatureTestCase;
use Tests\Helper\PlayerHelper;

/**
 * Class PlayerCommandsTest
 *
 * @package Tests\Feature\Console
 */
final class PlayerCommandsTest extends FeatureTestCase
{
    use DatabaseMigrations;
    use PlayerHelper;

    //region Tests

    /**
     * @param bool $withExistingName
     *
     * @return array
     */
    private function setUpCreatePlayerTest(bool $withExistingName = false): array
    {
        $name = $this->getFaker()->word;
        if ($withExistingName) {
            $this->createPlayerEntities(1, [PlayerModel::PROPERTY_NAME => $name]);
        }

        return [$name];
    }

    /**
     * @return void
     */
    public function testCreatePlayer(): void
    {
        [$name] = $this->setUpCreatePlayerTest();

        $this->artisan('player:create', ['name' => $name])
            ->expectsOutput(\sprintf('Player with name %s created', $name))
            ->assertExitCode(0);

        $this->assertNotEmpty($this->app->get(PlayerRepository::class)->findOneByName($name));
    }

    /**
     * @return void
     */
    public function testCreatePlayerWithExistingName(): void
    {
        [$name] = $this->setUpCreatePlayerTest(true);

        $this->artisan('player:create', ['name' => $name])
            ->expectsOutput(\sprintf('Player with name %s already exists', $name))
            ->assertExitCode(0);
    }

    //endregion
}
