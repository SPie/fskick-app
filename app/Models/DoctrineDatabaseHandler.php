<?php

namespace App\Models;

use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\EntityManager;
use Doctrine\ORM\Persisters\Entity\EntityPersister;

/**
 * Class DoctrineDatabaseHandler
 *
 * @package App\Models
 */
final class DoctrineDatabaseHandler implements DatabaseHandler
{
    /**
     * @var EntityManager
     */
    private EntityManager $entityManager;

    /**
     * @var string
     */
    private string $className;

    /**
     * DoctrineDatabaseHandler constructor.
     *
     * @param EntityManager $entityManager
     * @param string        $className
     */
    public function __construct(EntityManager $entityManager, string $className)
    {
        $this->entityManager = $entityManager;
        $this->className = $className;
    }

    /**
     * @return EntityManager
     */
    private function getEntityManager(): EntityManager
    {
        return $this->entityManager;
    }

    /**
     * @return string
     */
    private function getClassName(): string
    {
        return $this->className;
    }

    /**
     * @return EntityPersister
     */
    private function getEntityPersister(): EntityPersister
    {
        return $this->getEntityManager()->getUnitOfWork()->getEntityPersister($this->getClassName());
    }

    /**
     * @param int $id
     *
     * @return Model|object|null
     */
    public function find(int $id): ?Model
    {
        return $this->getEntityManager()->find($this->getClassName(), $id);
    }

    /**
     * @param array $criteria
     *
     * @return Model|object|null
     */
    public function load(array $criteria): ?Model
    {
        return $this->getEntityPersister()->load($criteria);
    }

    /**
     * @param array    $criteria
     * @param array    $orderBy
     * @param int|null $limit
     * @param int|null $offset
     *
     * @return Collection
     */
    public function loadAll(array $criteria = [], array $orderBy = [], int $limit = null, int $offset = null): Collection
    {
        return new ArrayCollection($this->getEntityPersister()->loadAll(
            $criteria,
            $orderBy,
            $limit,
            $offset
        ));
    }

    /**
     * @param Model $model
     * @param bool  $flush
     *
     * @return Model
     */
    public function save(Model $model, bool $flush = true): Model
    {
        $this->getEntityManager()->persist($model);

        if ($flush) {
            $this->flush();
        }

        return $model;
    }

    /**
     * @param Model $model
     * @param bool  $flush
     *
     * @return $this|DatabaseHandler
     */
    public function delete(Model $model, bool $flush = true): DatabaseHandler
    {
        $this->getEntityManager()->remove($model);

        if ($flush) {
            $this->flush();
        }

        return $this;
    }

    /**
     * @return $this|DatabaseHandler
     */
    public function flush(): DatabaseHandler
    {
        $this->getEntityManager()->flush();

        return $this;
    }
}
