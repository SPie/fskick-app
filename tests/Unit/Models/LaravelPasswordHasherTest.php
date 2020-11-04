<?php

namespace Tests\Unit\Models;

use App\Models\LaravelPasswordHasher;
use Illuminate\Hashing\HashManager;
use Mockery as m;
use Mockery\MockInterface;
use Tests\TestCase;

/**
 * Class LaravelPasswordHasherTest
 *
 * @package Tests\Unit\Models
 */
final class LaravelPasswordHasherTest extends TestCase
{
    //region Tests

    /**
     * @return void
     */
    public function testHash(): void
    {
        $password = $this->getFaker()->password;
        $hash = $this->getFaker()->sha256;
        $hashManager = $this->createHashManager();
        $this->mockHashManagerMake($hashManager, $hash, $password);

        $this->assertEquals($hash, $this->getLaravelPasswordHasher($hashManager)->hash($password));
    }

    /**
     * @param bool $validPassword
     *
     * @return array
     */
    private function setUpCheckTest(bool $validPassword = true): array
    {
        $password = $this->getFaker()->password;
        $hashedPassword = $this->getFaker()->sha256;
        $hashManager = $this->createHashManager();
        $this->mockHashManagerCheck($hashManager, $validPassword, $password, $hashedPassword);
        $passwordHasher = $this->getLaravelPasswordHasher($hashManager);

        return [$passwordHasher, $password, $hashedPassword];
    }

    /**
     * @return void
     */
    public function testCheck(): void
    {
        /** @var LaravelPasswordHasher $passwordHasher */
        [$passwordHasher, $password, $hashedPassword] = $this->setUpCheckTest();

        $this->assertTrue($passwordHasher->check($password, $hashedPassword));
    }

    /**
     * @return void
     */
    public function testCheckWithoutValidPassword(): void
    {
        /** @var LaravelPasswordHasher $passwordHasher */
        [$passwordHasher, $password, $hashedPassword] = $this->setUpCheckTest(false);

        $this->assertFalse($passwordHasher->check($password, $hashedPassword));
    }

    //endregion

    /**
     * @param HashManager|null $hashManager
     *
     * @return LaravelPasswordHasher
     */
    private function getLaravelPasswordHasher(HashManager $hashManager = null): LaravelPasswordHasher
    {
        return new LaravelPasswordHasher($hashManager ?: $this->createHashManager());
    }

    /**
     * @return HashManager|MockInterface
     */
    private function createHashManager(): HashManager
    {
        return m::spy(HashManager::class);
    }

    /**
     * @param HashManager|MockInterface $hashManager
     * @param string                    $hash
     * @param string                    $password
     *
     * @return $this
     */
    private function mockHashManagerMake(MockInterface $hashManager, string $hash, string $password): self
    {
        $hashManager
            ->shouldReceive('make')
            ->with($password)
            ->andReturn($hash);

        return $this;
    }

    /**
     * @param HashManager|MockInterface $hashManager
     * @param bool                      $valid
     * @param string                    $value
     * @param string                    $hashedValue
     *
     * @return $this
     */
    private function mockHashManagerCheck(MockInterface $hashManager, bool $valid, string $value, string $hashedValue): self
    {
        $hashManager
            ->shouldReceive('check')
            ->with($value, $hashedValue)
            ->andReturn($valid);

        return $this;
    }
}
