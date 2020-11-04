<?php

namespace App\Models;

/**
 * Interface Model
 *
 * @package App\Models
 */
interface Model
{
    const PROPERTY_ID = 'id';

    /**
     * @param int|null $id
     *
     * @return $this
     */
    public function setId(?int $id): self;

    /**
     * @return int|null
     */
    public function getId(): ?int;
}
