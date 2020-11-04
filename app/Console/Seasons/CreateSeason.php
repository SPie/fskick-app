<?php

namespace App\Console\Seasons;

use App\Models\Exceptions\ModelNotFoundException;
use App\Seasons\SeasonManager;
use Illuminate\Console\Command;

/**
 * Class CreateSeason
 *
 * @package App\Console\Seasons
 */
final class CreateSeason extends Command
{
    const ARGUMENT_NAME = 'name';

    /**
     * @var string
     */
    protected $signature = 'season:create {' . self::ARGUMENT_NAME . '}';

    /**
     * @var SeasonManager
     */
    private SeasonManager $seasonManager;

    /**
     * CreateSeason constructor.
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

        if ($this->seasonNameExists($name)) {
            $this->error(\sprintf('Season %s already exists', $name));

            return;
        }

        $season = $this->getSeasonManager()->createSeason($name);

        $this->line(\sprintf('Season %s created', $season->getName()));
    }

    /**
     * @param string $name
     *
     * @return bool
     */
    private function seasonNameExists(string $name): bool
    {
        try {
            $this->getSeasonManager()->getSeasonByName($name);
        } catch (ModelNotFoundException $e) {
            return false;
        }

        return true;
    }
}
