<?php

namespace App\Models;

/**
 * Interface UuidGenerator
 *
 * @package App\Models
 */
interface UuidGenerator
{
    /**
     * @return string
     */
    public function generate(): string;
}
