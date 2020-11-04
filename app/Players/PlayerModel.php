<?php

namespace App\Players;

use App\Models\Model;
use App\Models\Timestampable;

/**
 * Interface PlayerModel
 *
 * @package App\Players
 */
interface PlayerModel extends Model, Timestampable
{
    const PROPERTY_NAME = 'name';

    /**
     * @return string
     */
    public function getName(): string;
}
