<?php

namespace Tests\Helper;

/**
 * Trait ReflectionHelper
 *
 * @package Tests\Helper
 */
trait ReflectionHelper
{
    /**
     * @param $object
     *
     * @return \ReflectionObject
     */
    private function getReflectionObject($object)
    {
        return new \ReflectionObject($object);
    }

    /**
     * @param mixed  $object
     * @param string $propertyName
     *
     * @return mixed
     */
    private function getProtectedProperty($object, string $propertyName)
    {
        $property = $this->getReflectionObject($object)->getProperty($propertyName);
        $property->setAccessible(true);

        return $property->getValue($object);
    }
}
