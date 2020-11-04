<?php

namespace App\Console\Players;

use App\Models\Exceptions\ModelNotFoundException;
use App\Players\PlayerManager;
use Illuminate\Console\Command;

/**
 * Class CreatePlayer
 *
 * @package App\Console\Players
 */
final class CreatePlayer extends Command
{
    const ARGUMENT_NAME = 'name';

    /**
     * @var string
     */
    protected $signature = 'player:create {' . self::ARGUMENT_NAME . '}';

    /**
     * @var PlayerManager
     */
    private PlayerManager $playerManager;

    /**
     * CreatePlayer constructor.
     *
     * @param PlayerManager $playerManager
     */
    public function __construct(PlayerManager $playerManager)
    {
        parent::__construct();
        $this->playerManager = $playerManager;
    }

    /**
     * @return PlayerManager
     */
    private function getPlayerManager(): PlayerManager
    {
        return $this->playerManager;
    }

    /**
     * @return void
     */
    public function handle(): void
    {
        $name = $this->argument(self::ARGUMENT_NAME);
        if ($this->playerNameExists($name)) {
            $this->error(\sprintf('Player with name %s already exists', $name));
            return;
        }

        $player = $this->getPlayerManager()->createPlayer($name);

        $this->line(\sprintf('Player with name %s created', $player->getName()));
    }

    /**
     * @param string $name
     *
     * @return bool
     */
    private function playerNameExists(string $name): bool
    {
        try {
            $this->getPlayerManager()->getPlayerByName($name);
        } catch (ModelNotFoundException $e) {
            return false;
        }

        return true;
    }
}
