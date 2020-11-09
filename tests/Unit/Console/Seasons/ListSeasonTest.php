<?php

namespace Tests\Unit\Console\Seasons;

use App\Console\Seasons\ListSeason;
use App\Console\TableHelperFactory;
use App\Seasons\SeasonManager;
use App\Seasons\SeasonModel;
use Doctrine\Common\Collections\ArrayCollection;
use Tests\Helper\ConsoleHelper;
use Tests\Helper\ReflectionHelper;
use Tests\Helper\SeasonHelper;
use Tests\TestCase;

/**
 * Class ListSeasonTest
 *
 * @package Tests\Unit\Console\Seasons
 */
final class ListSeasonTest extends TestCase
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
            'season:list',
            $this->getProtectedProperty($this->getListSeason(), 'signature')
        );
    }

    /**
     * @param bool $withSeasons
     *
     * @return array
     */
    private function setUpHandleTest(bool $withSeasons = true): array
    {
        $name = $this->getFaker()->word;
        $active = $this->getFaker()->boolean;
        $season = $this->createSeasonModel();
        $this
            ->mockSeasonModelGetName($season, $name)
            ->mockSeasonModelIsActive($season, $active);
        $seasonManager = $this->createSeasonManager();
        $this->mockSeasonManagerGetSeasons($seasonManager, new ArrayCollection($withSeasons ? [$season] : []));
        $output = $this->createOutputStyle();
        $tableHelper = $this->createTableHelper();
        $tableHelperFactory = $this->createTableHelperFactory();
        $this->mockTableRendererRenderCreate($tableHelperFactory, $tableHelper, $output, ['Name', 'Active'], $withSeasons ? [[$name, $active]] : []);
        $command = $this->getListSeason($seasonManager, $tableHelperFactory);
        $this->setInputAndOutput($command, $this->createInput(), $output);

        return [$command, $tableHelper];
    }

    /**
     * @return void
     */
    public function testHandle(): void
    {
        /** @var ListSeason $command */
        [$command, $tableHelper] = $this->setUpHandleTest();

        $command->handle();

        $this->assertTableHelperRender($tableHelper);
    }

    /**
     * @return void
     */
    public function testHandleWithoutSeasons(): void
    {
        /** @var ListSeason $command */
        [$command, $tableHelper] = $this->setUpHandleTest(false);

        $command->handle();

        $this->assertTableHelperRender($tableHelper);
    }

    //endregion

    /**
     * @param SeasonManager|null      $seasonManager
     * @param TableHelperFactory|null $tableHelperFactory
     *
     * @return ListSeason
     */
    private function getListSeason(SeasonManager $seasonManager = null, TableHelperFactory $tableHelperFactory = null): ListSeason
    {
        return new ListSeason(
            $seasonManager ?: $this->createSeasonManager(),
            $tableHelperFactory ?: $this->createTableHelperFactory()
        );
    }
}
