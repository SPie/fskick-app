<?php

namespace Tests;

use Faker\Factory;
use Faker\Generator;

/**
 * Trait Faker
 *
 * @package Tests
 */
trait Faker
{
    /**
     * @var Generator
     */
    private Generator $faker;

    /**
     * @return Generator
     */
    protected function getFaker(): Generator
    {
        if (!isset($this->faker)) {
            $this->faker = Factory::create();
        }

        return $this->faker;
    }
}
