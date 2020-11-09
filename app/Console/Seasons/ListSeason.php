<?php

namespace App\Console\Seasons;

use App\Console\TableHelperFactory;
use App\Seasons\SeasonManager;
use App\Seasons\SeasonModel;
use Illuminate\Console\Command;

/**
 * Class ListSeason
 *
 * @package App\Console\Seasons
 */
final class ListSeason extends Command
{
    /**
     * @var string
     */
    protected $signature = 'season:list';

    /**
     * @var SeasonManager
     */
    private SeasonManager $seasonManager;

    /**
     * @var TableHelperFactory
     */
    private TableHelperFactory $tableHelperFactory;

    /**
     * ListSeason constructor.
     *
     * @param SeasonManager      $seasonManager
     * @param TableHelperFactory $tableHelperFactory
     */
    public function __construct(SeasonManager $seasonManager, TableHelperFactory $tableHelperFactory)
    {
        parent::__construct();

        $this->seasonManager = $seasonManager;
        $this->tableHelperFactory = $tableHelperFactory;
    }

    /**
     * @return SeasonManager
     */
    private function getSeasonManager(): SeasonManager
    {
        return $this->seasonManager;
    }

    /**
     * @return TableHelperFactory
     */
    private function getTableHelperFactory(): TableHelperFactory
    {
        return $this->tableHelperFactory;
    }

    /**
     * @return void
     */
    public function handle(): void
    {
        $tableHelper = $this->getTableHelperFactory()->create(
            $this->getOutput(),
            ['Name', 'Active'],
            $this->getSeasonManager()->getSeasons()
                ->map(fn (SeasonModel $season) => [$season->getName(), $season->isActive()])
                ->getValues()
        );

        $tableHelper->render();;
    }
}
