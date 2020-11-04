<?php

namespace Tests;

use Illuminate\Foundation\Testing\DatabaseMigrations as EloquentDatabaseMigrations;
use Illuminate\Foundation\Testing\TestCase as LaravelTestCase;
use LaravelDoctrine\Migrations\Testing\DatabaseMigrations;

/**
 * Class FeatureTestCase
 *
 * @package Tests
 */
abstract class FeatureTestCase extends LaravelTestCase
{
    use CreatesApplication;
    use Faker;

    /**
     * @return void
     */
    public function setUpTraits()
    {
        parent::setUpTraits();

        $uses = array_flip(class_uses_recursive(get_class($this)));

        // TODO mock queue service instead of email
        if (isset($uses[DatabaseMigrations::class]) && !isset($uses[EloquentDatabaseMigrations::class])) {
            $this->runDatabaseMigrations();
        }
    }
}
