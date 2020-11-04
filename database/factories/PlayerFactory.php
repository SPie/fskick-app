<?php

use App\Players\PlayerDoctrineModel;
use Faker\Generator as Faker;
use LaravelDoctrine\ORM\Testing\Factory;

/**
 * @var Factory $factory
 */

$factory->define(PlayerDoctrineModel::class, function (Faker $faker, array $attributes = []) {
    return [
        PlayerDoctrineModel::PROPERTY_NAME => $attributes[PlayerDoctrineModel::PROPERTY_NAME] ?? $faker->word,
    ];
});
