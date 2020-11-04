<?php

namespace App\Models;

use Doctrine\Common\Collections\Collection;

/**
 * Class AbstractDoctrineRepository
 *
 * @package App\Models
 */
abstract class AbstractDoctrineRepository implements Repository
{
    /**
     * @var DatabaseHandler
     */
    private DatabaseHandler $databaseHandler;

    /**
     * AbstractDoctrineRepository constructor.
     *
     * @param DatabaseHandler $databaseHandler
     */
    public function __construct(DatabaseHandler $databaseHandler)
    {
        $this->databaseHandler = $databaseHandler;
    }

    /**
     * @return DatabaseHandler
     */
    private function getDatabaseHandler(): DatabaseHandler
    {
        return $this->databaseHandler;
    }

    /**
     * @param int $id
     *
     * @return Model|null
     */
    public function find(int $id): ?Model
    {
        return $this->getDatabaseHandler()->find($id);
    }

    /**
     * @return Collection
     */
    public function findAll(): Collection
    {
        return $this->getDatabaseHandler()->loadAll();
    }

    /**
     * @param array $criteria
     *
     * @return Model|null
     */
    public function findOneBy(array $criteria): ?Model
    {
        return $this->getDatabaseHandler()->load($criteria);
    }

    /**
     * @param array      $criteria
     * @param array|null $orderBy
     * @param int|null   $limit
     * @param int|null   $offset
     *
     * @return Collection
     */
    public function findBy(array $criteria = [], array $orderBy = null, int $limit = null, int $offset = null): Collection
    {
        return $this->getDatabaseHandler()->loadAll($criteria, $orderBy, $limit, $offset);
    }

    /**
     * @param Model $model
     * @param bool  $flush
     *
     * @return Model
     */
    public function save(Model $model, bool $flush = true): Model
    {
        return $this->getDatabaseHandler()->save($model, $flush);
    }

    /**
     * @param Model $model
     * @param bool  $flush
     *
     * @return $this|Repository
     */
    public function delete(Model $model, bool $flush = true): Repository
    {
        $this->getDatabaseHandler()->delete($model, $flush);

        return $this;
    }

    /**
     * @return $this|Repository
     */
    public function flush(): Repository
    {
        $this->getDatabaseHandler()->flush();

        return $this;
    }
}
