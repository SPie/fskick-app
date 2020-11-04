<?php

namespace App\Models;

use Ramsey\Uuid\UuidFactoryInterface;

/**
 * Class RamseyUuidGenerator
 *
 * @package App\Models
 */
final class RamseyUuidGenerator implements UuidGenerator
{
    /**
     * @var UuidFactoryInterface
     */
    private UuidFactoryInterface $uuidFactory;

    /**
     * RamseyUuidGenerator constructor.
     *
     * @param UuidFactoryInterface $uuidFactory
     */
    public function __construct(UuidFactoryInterface $uuidFactory)
    {
        $this->uuidFactory = $uuidFactory;
    }

    /**
     * @return UuidFactoryInterface
     */
    private function getUuidFactory(): UuidFactoryInterface
    {
        return $this->uuidFactory;
    }

    /**
     * @return string
     */
    public function generate(): string
    {
        return $this->getUuidFactory()->uuid4()->toString();
    }
}
