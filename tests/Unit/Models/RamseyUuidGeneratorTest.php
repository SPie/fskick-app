<?php

namespace Tests\Unit\Models;

use App\Models\RamseyUuidGenerator;
use Mockery as m;
use Mockery\MockInterface;
use Ramsey\Uuid\UuidFactoryInterface;
use Ramsey\Uuid\UuidInterface;
use Tests\TestCase;

/**
 * Class RamseyUuidGeneratorTest
 *
 * @package Tests\Unit\Models
 */
final class RamseyUuidGeneratorTest extends TestCase
{
    //region Tests

    /**
     * @return void
     */
    public function testGenerate(): void
    {
        $uuidString = $this->getFaker()->uuid;
        $uuid = $this->createUuid();
        $this->mockUuidToString($uuid, $uuidString);
        $uuidFactory = $this->createUuidFactory();
        $this->mockUuidFactoryUuid4($uuidFactory, $uuid);
        $uuidGenerator = $this->getRamseyUuidGenerator($uuidFactory);

        $this->assertEquals($uuidString, $uuidGenerator->generate());
    }

    //endregion

    /**
     * @param UuidFactoryInterface $uuidFactory
     *
     * @return RamseyUuidGenerator
     */
    private function getRamseyUuidGenerator(UuidFactoryInterface $uuidFactory): RamseyUuidGenerator
    {
        return new RamseyUuidGenerator($uuidFactory);
    }

    /**
     * @return UuidFactoryInterface|MockInterface
     */
    private function createUuidFactory(): UuidFactoryInterface
    {
        return m::spy(UuidFactoryInterface::class);
    }

    /**
     * @param UuidFactoryInterface|MockInterface $uuidFactory
     * @param UuidInterface                      $uuid
     *
     * @return $this
     */
    private function mockUuidFactoryUuid4(MockInterface $uuidFactory, UuidInterface $uuid): self
    {
        $uuidFactory
            ->shouldReceive('uuid4')
            ->andReturn($uuid);

        return $this;
    }

    /**
     * @return UuidInterface|MockInterface
     */
    private function createUuid(): UuidInterface
    {
        return m::spy(UuidInterface::class);
    }

    /**
     * @param UuidInterface|MockInterface $uuid
     * @param string                      $uuidString
     *
     * @return $this
     */
    private function mockUuidToString(MockInterface $uuid, string $uuidString): self
    {
        $uuid
            ->shouldReceive('toString')
            ->andReturn($uuidString);

        return $this;
    }
}
