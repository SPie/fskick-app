<?php

namespace Tests\Unit\Console\Seasons;

use App\Console\Seasons\ActivateSeason;
use App\Models\Exceptions\ModelNotFoundException;
use App\Seasons\SeasonManager;
use App\Seasons\SeasonModel;
use Tests\Helper\ConsoleHelper;
use Tests\Helper\ReflectionHelper;
use Tests\Helper\SeasonHelper;
use Tests\TestCase;

/**
 * Class ActivateSeasonTest
 *
 * @package Tests\Unit\Console\Seasons
 */
final class ActivateSeasonTest extends TestCase
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
            'season:activate {name}',
            $this->getProtectedProperty($this->getActivateSeason(), 'signature')
        );
    }

    /**
     * @param bool $withSeason
     *
     * @return array
     */
    private function setUpHandleTest(bool $withSeason = true): array
    {
        $name = $this->getFaker()->word;
        $input = $this->createInput();
        $this->mockInputGetArgument($input, $name, 'name');
        $season = $this->createSeasonModel();
        $activatedSeason = $this->createSeasonModel();
        $this->mockSeasonModelGetName($activatedSeason, $name);
        $seasonManager = $this->createSeasonManager();
        $this
            ->mockSeasonManagerGetSeasonByName($seasonManager, $withSeason ? $season : new ModelNotFoundException(), $name)
            ->mockSeasonManagerActivateSeason($seasonManager, $activatedSeason, $season);
        $output = $this->createOutputStyle();
        $command = $this->getActivateSeason($seasonManager);
        $this->setInputAndOutput($command, $input, $output);

        return [$command, $activatedSeason, $output];
    }

    /**
     * @return void
     */
    public function testHandle(): void
    {
        /**
         * @var ActivateSeason $command
         * @var SeasonModel    $activatedSeason
         */
        [$command, $activatedSeason, $output] = $this->setUpHandleTest();

        $command->handle();

        $this->assertOutputStyleWriteln($output, \sprintf('Season %s is active now', $activatedSeason->getName()), 32);
    }

    /**
     * @return void
     */
    public function testHandleWithoutSeason(): void
    {
        /**
         * @var ActivateSeason $command
         * @var SeasonModel    $activatedSeason
         */
        [$command, $activatedSeason, $output] = $this->setUpHandleTest(false);

        $command->handle();

        $this->assertOutputStyleWriteln(
            $output,
            \sprintf('<error>Season with name %s doesn\'t exist</error>', $activatedSeason->getName()),
            32
        );
    }

    //endregion

    /**
     * @param SeasonManager|null $seasonManager
     *
     * @return ActivateSeason
     */
    private function getActivateSeason(SeasonManager $seasonManager = null): ActivateSeason
    {
        return new ActivateSeason($seasonManager ?: $this->createSeasonManager());
    }
}
