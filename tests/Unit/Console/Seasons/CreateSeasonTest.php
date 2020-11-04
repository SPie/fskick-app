<?php

namespace Tests\Unit\Console\Seasons;

use App\Console\Seasons\CreateSeason;
use App\Models\Exceptions\ModelNotFoundException;
use App\Seasons\SeasonManager;
use App\Seasons\SeasonModel;
use Tests\Helper\ConsoleHelper;
use Tests\Helper\ReflectionHelper;
use Tests\Helper\SeasonHelper;
use Tests\TestCase;

/**
 * Class CreateSeasonTest
 *
 * @package Tests\Unit\Console\Seasons
 */
final class CreateSeasonTest extends TestCase
{
    use ConsoleHelper;
    use ReflectionHelper;
    use SeasonHelper;

    //region Tests

    /**
     * @return void
     */
    public function testSignature(): void
    {
        $this->assertEquals(
            'season:create {name}',
            $this->getProtectedProperty($this->getCreateSeason(), 'signature')
        );
    }

    /**
     * @param bool $withExistingName
     *
     * @return array
     */
    private function setUpHandleTest(bool $withExistingName = false): array
    {
        $name = $this->getFaker()->word;
        $input = $this->createInput();
        $this->mockInputGetArgument($input, $name, 'name');
        $season = $this->createSeasonModel();
        $this->mockSeasonModelGetName($season, $name);
        $output = $this->createOutputStyle();
        $seasonManager = $this->createSeasonManager();
        $this
            ->mockSeasonManagerCreateSeason($seasonManager, $season, $name)
            ->mockSeasonManagerGetSeasonByName($seasonManager, $withExistingName ? $this->createSeasonModel() : new ModelNotFoundException(), $name);
        $command = $this->getCreateSeason($seasonManager);
        $this->setInputAndOutput($command, $input, $output);

        return [$command, $season, $output, $seasonManager];
    }

    /**
     * @return void
     */
    public function testHandle(): void
    {
        /**
         * @var CreateSeason $command
         * @var SeasonModel  $season
         */
        [$command, $season, $output, $seasonManager] = $this->setUpHandleTest();

        $command->handle();

        $this->assertOutputStyleWriteln($output, \sprintf('Season %s created', $season->getName()), 32);
        $this->assertSeasonManagerCreateSeason($seasonManager, $season->getName());
    }

    /**
     * @return void
     */
    public function testHandleWithExistingName(): void
    {
        /**
         * @var CreateSeason $command
         * @var SeasonModel  $season
         */
        [$command, $season, $output] = $this->setUpHandleTest(true);

        $command->handle();

        $this->assertOutputStyleWriteln($output, \sprintf('<error>Season %s already exists</error>', $season->getName()), 32);
    }

    //endregion

    /**
     * @param SeasonManager|null $seasonManager
     *
     * @return CreateSeason
     */
    private function getCreateSeason(SeasonManager $seasonManager = null): CreateSeason
    {
        return new CreateSeason($seasonManager ?: $this->createSeasonManager());
    }
}
