<?php

use App\Seasons\SeasonDoctrineModel;
use Faker\Generator as Faker;
use LaravelDoctrine\ORM\Testing\Factory;

/**
 * @var Factory $factory
 */

$factory->define(SeasonDoctrineModel::class, function (Faker $faker, array $attributes = []) {
    return [
        SeasonDoctrineModel::PROPERTY_NAME   => $attributes[SeasonDoctrineModel::PROPERTY_NAME] ?? $faker->word,
        SeasonDoctrineModel::PROPERTY_ACTIVE => $attributes[SeasonDoctrineModel::PROPERTY_ACTIVE] ?? $faker->boolean,
    ];
});
