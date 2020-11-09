<?php

namespace Tests\Feature\Console;

use App\Seasons\SeasonModel;
use App\Seasons\SeasonRepository;
use LaravelDoctrine\Migrations\Testing\DatabaseMigrations;
use Tests\FeatureTestCase;
use Tests\Helper\SeasonHelper;

/**
 * Class SeasonCommandsTest
 *
 * @package Tests\Feature\Console
 */
final class SeasonCommandsTest extends FeatureTestCase
{
    use DatabaseMigrations;
    use SeasonHelper;

    //region Tests

    /**
     * @param bool $withExistingName
     *
     * @return array
     */
    private function setUpCreateSeasonTest(bool $withExistingName = false): array
    {
        $name = $this->getFaker()->word;
        if ($withExistingName) {
            $this->createSeasonEntities(1, [SeasonModel::PROPERTY_NAME => $name]);
        }

        return [$name];
    }

    /**
     * @return void
     */
    public function testCreateSeason(): void
    {
        [$name] = $this->setUpCreateSeasonTest();

        $this->artisan('season:create', ['name' => $name])
            ->expectsOutput(\sprintf('Season %s created', $name))
            ->assertExitCode(0);

        $this->assertNotEmpty($this->app->get(SeasonRepository::class)->findOneByName($name));
    }

    /**
     * @return void
     */
    public function testCreateSeasonWithExistingName(): void
    {
        [$name] = $this->setUpCreateSeasonTest(true);

        $this->artisan('season:create', ['name' => $name])
            ->expectsOutput(\sprintf('Season %s already exists', $name))
            ->assertExitCode(0);
    }

    /**
     * @return void
     */
    public function testListSeasons(): void
    {
        $season = $this->createSeasonEntities()->first();

        $this->artisan('season:list')
            ->expectsTable(['Name', 'Active'], [[$season->getName(), $season->isActive()]])
            ->assertExitCode(0);

        $this->expectedTables = [];
    }

    //endregion
}
