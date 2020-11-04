<?php

namespace App\Providers;

use App\Models\DatabaseHandler;
use App\Models\DoctrineDatabaseHandler;
use App\Models\LaravelPasswordHasher;
use App\Models\PasswordHasher;
use App\Models\RamseyUuidGenerator;
use App\Models\UuidGenerator;
use App\Players\PlayerDoctrineModel;
use App\Players\PlayerDoctrineModelFactory;
use App\Players\PlayerDoctrineRepository;
use App\Players\PlayerModel;
use App\Players\PlayerModelFactory;
use App\Players\PlayerRepository;
use App\Seasons\SeasonDoctrineModel;
use App\Seasons\SeasonDoctrineModelFactory;
use App\Seasons\SeasonDoctrineRepository;
use App\Seasons\SeasonModel;
use App\Seasons\SeasonModelFactory;
use App\Seasons\SeasonRepository;
use Doctrine\ORM\EntityManager;
use Illuminate\Contracts\Container\Container;
use Illuminate\Hashing\HashManager;
use Illuminate\Support\ServiceProvider;
use Ramsey\Uuid\UuidFactory;

/**
 * Class ModelServiceProvider
 *
 * @package App\Providers
 */
final class ModelServiceProvider extends ServiceProvider
{

    /**
     * @return void
     */
    public function register()
    {
        $this
            ->bindModels()
            ->bindModelFactories()
            ->bindDatabaseHandler()
            ->bindRepositories()
            ->bindUuidGenerator()
            ->bindPasswordHasher();
    }

    /**
     * @return $this
     */
    private function bindModels(): self
    {
        $this->app->bind(PlayerModel::class, PlayerDoctrineModel::class);
        $this->app->bind(SeasonModel::class, SeasonDoctrineModel::class);

        return $this;
    }

    /**
     * @return $this
     */
    private function bindModelFactories(): self
    {
        $this->app->singleton(PlayerModelFactory::class, PlayerDoctrineModelFactory::class);
        $this->app->singleton(SeasonModelFactory::class, SeasonDoctrineModelFactory::class);

        return $this;
    }

    /**
     * @return $this
     */
    private function bindDatabaseHandler(): self
    {
        $this->app->bind(
            DatabaseHandler::class,
            fn (Container $app, array $parameters) => new DoctrineDatabaseHandler($parameters[0], $parameters[1])
        );

        return $this;
    }

    /**
     * @param string $className
     *
     * @return DatabaseHandler
     */
    private function makeDatabaseHandler(string $className): DatabaseHandler
    {
        return $this->app->make(DatabaseHandler::class, [$this->app->get(EntityManager::class), $className]);
    }

    /**
     * @return $this
     */
    private function bindRepositories(): self
    {
        $this->app->singleton(
            PlayerRepository::class,
            fn () => new PlayerDoctrineRepository($this->makeDatabaseHandler(PlayerDoctrineModel::class))
        );
        $this->app->singleton(
            SeasonRepository::class,
            fn () => new SeasonDoctrineRepository($this->makeDatabaseHandler(SeasonDoctrineModel::class))
        );

        return $this;
    }

    /**
     * @return $this
     */
    private function bindUuidGenerator(): self
    {
        $this->app->singleton(UuidGenerator::class, fn () => new RamseyUuidGenerator(new UuidFactory()));

        return $this;
    }

    /**
     * @return $this
     */
    private function bindPasswordHasher(): self
    {
        $this->app->singleton(
            PasswordHasher::class,
            fn () => new LaravelPasswordHasher($this->app->get(HashManager::class))
        );

        return $this;
    }
}
