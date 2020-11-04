<?php

namespace Tests\Unit\Console\Players;

use App\Console\Players\CreatePlayer;
use App\Models\Exceptions\ModelNotFoundException;
use App\Players\PlayerManager;
use App\Players\PlayerModel;
use Tests\Helper\ConsoleHelper;
use Tests\Helper\PlayerHelper;
use Tests\Helper\ReflectionHelper;
use Tests\TestCase;

/**
 * Class CreatePlayerTest
 *
 * @package App\Console\Players
 */
final class CreatePlayerTest extends TestCase
{
    use ConsoleHelper;
    use PlayerHelper;
    use ReflectionHelper;

    //region Tests

    /**
     * @return void
     */
    public function testSignature(): void
    {
        $this->assertEquals(
            'player:create {name}',
            $this->getProtectedProperty($this->getCreatePlayer(), 'signature')
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
        $player = $this->createPlayerModel();
        $this->mockPlayerModelGetName($player, $name);
        $playerManager = $this->createPlayerManager();
        $this
            ->mockPlayerManagerCreatePlayer($playerManager, $player, $name)
            ->mockPlayerManagerGetPlayerByName($playerManager, $withExistingName ? $player : new ModelNotFoundException(), $name);
        $output = $this->createOutputStyle();
        $command = $this->getCreatePlayer($playerManager);
        $this->setInputAndOutput($command, $input, $output);

        return [$command, $player, $output, $playerManager];
    }

    /**
     * @return void
     */
    public function testHandle(): void
    {
        /**
         * @var CreatePlayer $command
         * @var PlayerModel  $player
         */
        [$command, $player, $output, $playerManager] = $this->setUpHandleTest();

        $command->handle();

        $this->assertOutputStyleWriteln($output, \sprintf('Player with name %s created', $player->getName()), 32);
        $this->assertPlayerManagerCreatePlayer($playerManager, $player->getName());
    }

    /**
     * @return void
     */
    public function testHandleWithExistingName(): void
    {
        /**
         * @var CreatePlayer $command
         * @var PlayerModel  $player
         */
        [$command, $player, $output] = $this->setUpHandleTest(true);

        $command->handle();

        $this->assertOutputStyleWriteln(
            $output,
            \sprintf('<error>Player with name %s already exists</error>', $player->getName()),
            32
        );
    }

    //endregion

    /**
     * @param PlayerManager|null $playerManager
     *
     * @return CreatePlayer
     */
    private function getCreatePlayer(PlayerManager $playerManager = null): CreatePlayer
    {
        return new CreatePlayer($playerManager ?: $this->createPlayerManager());
    }
}
