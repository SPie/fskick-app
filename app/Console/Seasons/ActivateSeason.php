<?php

namespace App\Console\Seasons;

use App\Models\Exceptions\ModelNotFoundException;
use App\Seasons\SeasonManager;
use App\Seasons\SeasonModel;
use Illuminate\Console\Command;

/**
 * Class ActivateSeason
 *
 * @package App\Console\Seasons
 */
final class ActivateSeason extends Command
{
    const ARGUMENT_NAME = 'name';

    protected $signature = 'season:activate {' . self::ARGUMENT_NAME . '}';

    /**
     * @var SeasonManager
     */
    private SeasonManager $seasonManager;

    /**
     * ActivateSeason constructor.
     *
     * @param SeasonManager $seasonManager
     */
    public function __construct(SeasonManager $seasonManager)
    {
        parent::__construct();

        $this->seasonManager = $seasonManager;
    }

    /**
     * @return SeasonManager
     */
    private function getSeasonManager(): SeasonManager
    {
        return $this->seasonManager;
    }

    /**
     * @return void
     */
    public function handle(): void
    {
        $name = $this->argument(self::ARGUMENT_NAME);
        $season = $this->getSeason($name);
        if (!$season) {
            $this->error(\sprintf('Season with name %s doesn\'t exist', $name));
            return;

        }

        $season = $this->getSeasonManager()->activateSeason($season);

        $this->line(\sprintf('Season %s is active now', $season->getName()));
    }

    /**
     * @param string $name
     *
     * @return SeasonModel|null
     */
    private function getSeason(string $name): ?SeasonModel
    {
        try {
            return $this->getSeasonManager()->getSeasonByName($name);
        } catch (ModelNotFoundException $e) {
            return null;
        }
    }
}
