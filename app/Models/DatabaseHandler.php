<?php

namespace App\Models;

use Doctrine\Common\Collections\Collection;

/**
 * Interface DatabaseHandler
 *
 * @package App\Models
 */
interface DatabaseHandler
{
    /**
     * @param int $id
     *
     * @return Model|null
     */
    public function find(int $id): ?Model;

    /**
     * @param array $criteria
     *
     * @return Model|null
     */
    public function load(array $criteria): ?Model;

    /**
     * @param array    $criteria
     * @param array    $orderBy
     * @param int|null $limit
     * @param int|null $offset
     *
     * @return Collection
     */
    public function loadAll(array $criteria = [], array $orderBy = [], int $limit = null, int $offset = null): Collection;

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
