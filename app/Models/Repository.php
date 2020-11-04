<?php

namespace App\Models;

use Doctrine\Common\Collections\Collection;

/**
 * Interface Repository
 *
 * @package App\Models
 */
interface Repository
{
    /**
     * @param int $id
     *
     * @return Model|null
     */
    public function find(int $id): ?Model;

    /**
     * @return Collection
     */
    public function findAll(): Collection;

    /**
     * @param array $criteria
     *
     * @return Model|null
     */
    public function findOneBy(array $criteria): ?Model;

    /**
     * @param array      $criteria
     * @param array|null $orderBy
     * @param int|null   $limit
     * @param int|null   $offset
     *
     * @return Collection
     */
    public function findBy(array $criteria = [], array $orderBy = null, int $limit = null, int $offset = null): Collection;

    /**
     * @param Model $model
     * @param bool  $flush
     *
     * @return Model
     */
    public function save(Model $model, bool $flush = true): Model;

    /**
     * @param Model $model
     * @param bool  $flush
     *
     * @return $this
     */
    public function delete(Model $model, bool $flush = true): self;

    /**
     * @return $this
     */
    public function flush(): self;
}
