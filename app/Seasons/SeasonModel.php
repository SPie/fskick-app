<?php

namespace App\Seasons;

use App\Models\Model;
use App\Models\Timestampable;

/**
 * Interface SeasonModel
 *
 * @package App\Seasons
 */
interface SeasonModel extends Model, Timestampable
{
    const PROPERTY_NAME   = 'name';
    const PROPERTY_ACTIVE = 'active';

    /**
     * @return string
     */
    public function getName(): string;

    /**
     * @param bool $active
     *
     * @return $this
     */
    public function setActive(bool $active): self;

    /**
     * @return bool
     */
    public function isActive(): bool;
}
